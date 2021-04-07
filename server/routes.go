package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
)

func uploadFileRoute(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get uploaded file.")
	}

	owner := "solaris" //ctx.Locals("user").(string)
	uploadedFile, uploadErr := uploadFile(owner, fileHeader)

	if uploadErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, uploadErr.Error())
	}

	b, jsonErr := json.Marshal(ImageResponse{
		Url:     fmt.Sprintf("%v/%v", cdnConfig.CdnEndpoint, uploadedFile),
		Success: true,
	})

	if jsonErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, jsonErr.Error())
	}

	return ctx.Send(b)
}

func getFileInfoRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	rec, err := getFileRecord(id)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "An internal error occurred while retriving file.")
	}

	file, _ := json.Marshal(rec)
	return ctx.Send(file)
}

func getFileRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	rec, err := getFileRecord(id)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "An internal error occurred while retriving file.")
	}

	imageURL := fmt.Sprintf("%s/%s%s", cdnConfig.SpacesConfig.SpacesUrl, rec.Name, rec.Ext)

	res, err := http.Head(imageURL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fiber.NewError(res.StatusCode, "An error occurred redirecting to the image")
	}

	return ctx.Redirect(imageURL)
}

func getFilesRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	folder, files, err := getFolderRecord(id)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "An internal error occurred while retriving folder info.")
	}

	folderMap := structs.Map(folder)
	delete(folderMap, "files")
	folderMap["files"] = files

	folderJSON, _ := json.Marshal(folderMap)
	return ctx.Send(folderJSON)
}
