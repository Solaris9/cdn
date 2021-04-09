package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func getOGEmbedRoute(ctx *fiber.Ctx) error {
	id, ext := ctx.Params("id"), ctx.Params("ext")
	imageURL := fmt.Sprintf("%v/%v.%v", cdnConfig.SpacesConfig.SpacesUrl, id, ext)

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

		return ctx.Send(toJSON(jsonObj))
	} else {
		return ctx.Redirect(imageURL, fiber.StatusMovedPermanently)
	}
}

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
	id, ext := ctx.Params("id"), ctx.Params("ext")
	imageURL := fmt.Sprintf("%v/%v.%v", cdnConfig.SpacesConfig.SpacesUrl, id, ext)

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

func createFolderRoute(ctx *fiber.Ctx) error {
	body := new(FolderPostRequest)

	if err := ctx.BodyParser(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if body.Name == "" {
		respErr := NewResponse(fiber.StatusBadRequest, "Folder name required.")
		return ctx.Send(respErr.ToJSON())
	}

	folder, respErr := NewFolder(body.Name, ctx.Locals("user").(string))
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	return ctx.Send(folder.ToJSON())
}

func getFolderRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	folder, respErr := FolderFor(id, true)

	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	return ctx.Send(folder.ToJSON())
}

func updateFolderRoute(ctx *fiber.Ctx) error {
	body := new(FolderPatchRequest)
	if err := ctx.BodyParser(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id := ctx.Params("id")
	folder, respErr := FolderFor(id, false)
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	respErr = folder.CheckOwner(ctx.Locals("user").(string))
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
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
			return ctx.Send(respErr.ToJSON())
		}
	}

	return ctx.Send(folder.ToJSON())
}

func deleteFolderRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	folder, respErr := FolderFor(id, true)
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	respErr = folder.CheckOwner(ctx.Locals("user").(string))
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	respErr = folder.Delete()
	if respErr != nil {
		return ctx.Send(respErr.ToJSON())
	}

	return ctx.Send(folder.ToJSON())
}
