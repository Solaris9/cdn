package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber/v2"
)

func getOGEmbedRoute(ctx *fiber.Ctx) error {
	file := ctx.Params("file")
	imageURL := fmt.Sprintf("%v/%v", cdnConfig.SpacesConfig.SpacesUrl, file)

	res, headErr := http.Head(imageURL)
	if headErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, headErr.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fiber.NewError(res.StatusCode, "An error occurred redirecting to the image.")
	}

	if ctx.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
		ctx.Type("json", "utf-8")
		contentLengthHeader := res.Header.Values("content-length")[0]
		contentTypeHeader := res.Header.Values("content-type")[0]
		contentLength, parseErr := strconv.ParseInt(contentLengthHeader, 0, 64)
		if parseErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, parseErr.Error())
		}

		var objType string
		objAuthor := fmt.Sprintf("%v | %v", getFileSize(contentLength), contentTypeHeader)
		objProvider := res.Header.Values("last-modified")[0]

		if strings.HasPrefix(contentTypeHeader, "image") {
			objType = "photo"
		} else if strings.HasPrefix(contentTypeHeader, "video") {
			objType = "video"
		} else {
			objType = "link"
		}

		jsonObj := Embed{
			Type:         objType,
			AuthorName:   objAuthor,
			ProviderName: objProvider,
		}

		return ctx.JSON(jsonObj)
	} else {
		return ctx.Redirect(imageURL, fiber.StatusMovedPermanently)
	}
}

func uploadFileRoute(ctx *fiber.Ctx) error {
	s, err := session.NewSession(cdnS3Config)
	if err != nil {
		respErr := NewResponseByError(fiber.StatusBadRequest, err)
		return ctx.JSON(respErr)
	}

	file, respErr := UploadFile(s, ctx)
	if respErr != nil {
		return ctx.JSON(respErr)
	}

	url := fmt.Sprintf("%v/%v", cdnConfig.SpacesConfig.SpacesUrl, file)

	resp := ImageResult{
		Code:    200,
		Url:     url,
		Success: true,
	}

	return ctx.JSON(resp)
}

func getFileRoute(ctx *fiber.Ctx) error {
	file := ctx.Params("file")
	imageURL := fmt.Sprintf("%v/%v", cdnConfig.SpacesConfig.SpacesUrl, file)

	res, err := http.Head(imageURL)
	if err != nil {
		respErr := NewResponseByError(fiber.StatusInternalServerError, err)
		return ctx.JSON(respErr)
	}

	if res.StatusCode != http.StatusOK {
		respErr := NewResponse(res.StatusCode, "An error occurred redirecting to the image")
		return ctx.JSON(respErr)
	}

	return ctx.Redirect(imageURL)
}

func getFilesRoute(ctx *fiber.Ctx) error {
	s, err := session.NewSession(cdnS3Config)
	if err != nil {
		respErr := NewResponseByError(fiber.StatusBadRequest, err)
		return ctx.JSON(respErr)
	}

	objects, objectsErr := GetFiles(s)
	if objectsErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, objectsErr.Error())
	}

	data := fiber.Map{
		"files":  objects,
		"length": len(objects),
	}

	return ctx.JSON(data)
}

func deleteFileRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	s, err := session.NewSession(cdnS3Config)
	if err != nil {
		respErr := NewResponseByError(fiber.StatusBadRequest, err)
		return ctx.JSON(respErr)
	}

	respErr := CheckOwner(s, id, ctx.Locals("user").(string))
	if respErr != nil {
		return ctx.JSON(respErr)
	}

	respErr = DeleteFile(s, id)
	if respErr != nil {
		return ctx.JSON(respErr)
	}

	return ctx.JSON(fiber.Map{
		"id":      id,
		"success": true,
		"code":    200,
	})
}

func createFolderRoute(ctx *fiber.Ctx) error {
	ctx.Set("content-type", "application/json")
	body := new(FolderPostRequest)

	if err := ctx.BodyParser(body); err != nil {
		respErr := NewResponseByError(fiber.StatusBadRequest, err)
		return ctx.JSON(respErr)
	}

	if body.Name == "" {
		respErr := NewResponse(fiber.StatusBadRequest, "Folder name required.")
		return ctx.JSON(respErr)
	}

	folder, respErr := NewFolder(body.Name, ctx.Locals("user").(string))
	if respErr != nil {
		return ctx.JSON(respErr)
	}

	return ctx.JSON(&FolderResult{
		CreateTime: folder.CreateTime,
		UpdateTime: folder.UpdateTime,
		ID:         folder.Data.ID,
		Name:       folder.Data.Name,
	})
}

func getFoldersRoute(ctx *fiber.Ctx) error {
	ctx.Set("content-type", "application/json")

	firebaseCtx := context.Background()
	docs, err := cdnFirestore.Collection("folders").Documents(firebaseCtx).GetAll()
	if err != nil {
		respErr := NewResponseByError(fiber.StatusInternalServerError, err)
		return ctx.JSON(respErr)
	}

	folders := make([]*FoldersResult, len(docs))
	for i, doc := range docs {
		data := new(FolderData)
		doc.DataTo(data)

		folders[i] = &FoldersResult{
			CreateTime: doc.CreateTime,
			UpdateTime: doc.UpdateTime,
			ID:         data.ID,
			Name:       data.Name,
			Size:       len(data.Files),
		}
	}

	return ctx.JSON(folders)
}

func getFolderRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	folder, respErr := FolderFor(id)

	if respErr != nil {
		return ctx.JSON(respErr)
	}

	s, err := session.NewSession(cdnS3Config)
	if err != nil {
		respErr := NewResponseByError(fiber.StatusInternalServerError, err)
		return ctx.JSON(respErr)
	}

	files, err := GetFilesByKeys(s, folder.Data.Files)
	if err != nil {
		respErr := NewResponseByError(fiber.StatusInternalServerError, err)
		return ctx.JSON(respErr)
	}

	return ctx.JSON(&FolderResult{
		CreateTime: folder.CreateTime,
		UpdateTime: folder.UpdateTime,
		ID:         folder.Data.ID,
		Name:       folder.Data.Name,
		Files:      files,
	})
}

func updateFolderRoute(ctx *fiber.Ctx) error {
	body := new(FolderPatchRequest)

	if err := ctx.BodyParser(body); err != nil {
		respErr := NewResponseByError(fiber.StatusBadRequest, err)
		return ctx.JSON(respErr)
	}

	id := ctx.Params("id")
	folder, respErr := FolderFor(id)
	if respErr != nil {
		return ctx.JSON(respErr)
	}

	respErr = folder.CheckOwner(ctx.Locals("user").(string))
	if respErr != nil {
		return ctx.JSON(respErr)
	}

	if body.Name != "" {
		folder.Data.Name = body.Name
	}

	if body.Add != nil {
		folder.AddFiles(body.Add, false)
	}

	if body.Remove != nil {
		folder.RemoveFiles(body.Remove)
	}

	if folder.IsChanged() {
		if respErr := folder.Save(); respErr != nil {
			return ctx.JSON(respErr)
		}
	}

	return ctx.JSON(&FolderResult{
		CreateTime: folder.CreateTime,
		UpdateTime: folder.UpdateTime,
		ID:         folder.Data.ID,
		Name:       folder.Data.Name,
	})
}

func deleteFolderRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	folder, respErr := FolderFor(id)
	if respErr != nil {
		return ctx.JSON(respErr)
	}

	respErr = folder.CheckOwner(ctx.Locals("user").(string))
	if respErr != nil {
		return ctx.JSON(respErr)
	}

	respErr = folder.Delete()
	if respErr != nil {
		return ctx.JSON(respErr)
	}

	return ctx.JSON(&FolderResult{
		CreateTime: folder.CreateTime,
		UpdateTime: folder.UpdateTime,
		ID:         folder.Data.ID,
		Name:       folder.Data.Name,
	})
}
