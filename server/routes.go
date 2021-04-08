package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func uploadFileRoute(ctx *fiber.Ctx) error {
	file, respErr := UploadFile(ctx)
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	resp := ImageResponse{
		Code:    200,
		Url:     file.URL(),
		Success: true,
	}

	return ctx.Send(toJSON(resp))
}

func getFileRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	file, respErr := FileFor(id)

	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	imageURL := file.ImageURL()

	res, err := http.Head(imageURL)
	if err != nil {
		respErr := NewResponseByError(fiber.StatusInternalServerError, err)
		return ctx.Send(respErr.ToJSON())
	}

	if res.StatusCode != http.StatusOK {
		respErr := NewResponse(res.StatusCode, "An error occurred redirecting to the image")
		return ctx.Send(respErr.ToJSON())
	}

	return ctx.Redirect(imageURL)
}

func getFileInfoRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	file, respErr := FileFor(id)

	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	return ctx.Send(file.ToJSON())
}

func deleteFileRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	file, respErr := FileFor(id)
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	respErr = file.CheckOwner(ctx.Locals("user").(string))
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	respErr = file.Delete()
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	return ctx.Send(toJSON(fiber.Map{
		"id":      id,
		"success": true,
		"code":    200,
	}))
}

// func createFolderRoute(ctx *fiber.Ctx) error {
// 	owner := "solaris" //ctx.Locals("user").(string)
// 	body := new(FolderPostRequest)

// 	if err := ctx.BodyParser(body); err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, err.Error())
// 	}

// 	folder, err := createFolderRecord(body.Name, owner)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "An internal error occurred while retriving folder info.")
// 	}

// 	folderJSON, _ := json.Marshal(fiber.Map{
// 		"name": folder.Name,
// 		"id":   folder.ID,
// 	})

// 	return ctx.Send(folderJSON)
// }

// func getFolderRoute(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")
// 	folder, files, err := getFolderRecord(id)

// 	if err != nil {
// 		if status.Code(err) == codes.NotFound {
// 			return fiber.NewError(fiber.StatusNotFound, "This folder does not exist.")
// 		}

// 		return fiber.NewError(fiber.StatusInternalServerError, "An internal error occurred while retriving folder info.")
// 	}

// 	folderMap := structs.Map(folder)
// 	delete(folderMap, "files")
// 	folderMap["files"] = files

// 	folderJSON, _ := json.Marshal(folderMap)
// 	return ctx.Send(folderJSON)
// }

// func updateFolderRoute(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")
// 	owner := "solaris" //ctx.Locals("user").(string)
// 	body := new(FolderPatchRequest)

// 	if err := ctx.BodyParser(body); err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, err.Error())
// 	}

// 	folder, err := getFolderRecord(id)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "An internal error occurred while retriving folder info.")
// 	}

// 	success, err := addFilesToFolderRecord(id, owner, body.Files)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "An internal error occurred while retriving folder info.")
// 	}

// 	folderJSON, _ := json.Marshal(fiber.Map{
// 		"name": folder.Name,
// 		"id":   folder.ID,
// 	})

// 	return ctx.Send(folderJSON)
// }
