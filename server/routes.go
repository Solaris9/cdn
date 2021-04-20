package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber/v2"
)

func verifyAuthRoute(ctx *fiber.Ctx) error {
	body := new(TokenResponse)

	if err := ctx.BodyParser(body); err != nil {
		respErr := NewResponseByError(fiber.StatusBadRequest, err)
		respErr.SetSuccess(false)
		return ctx.JSON(respErr)
	}

	if body.Token == "" {
		return ctx.JSON(&TokenRequest{
			Success: false,
			Message: "No authorization token provided.",
		})
	}

	if body.Token != cdnConfig.Authorization {
		return ctx.JSON(&TokenRequest{
			Success: false,
			Message: "Invalid authorization token provided.",
		})
	}

	return ctx.JSON(&TokenRequest{
		Success: true,
	})
}

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

	url := fmt.Sprintf("%v/%v", cdnConfig.CdnEndpoint, file)

	resp := ImageResult{
		Code:    200,
		Url:     url,
		Success: true,
	}

	return ctx.JSON(resp)
}

func getFileRoute(ctx *fiber.Ctx) error {
	key := ctx.Params("id")

	fmt.Println("Endpoint Hit: getImage")

	queries := new(ImageResponseQuery)

	if queryErr := ctx.QueryParser(queries); queryErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, queryErr.Error())
	}

	imageURL := fmt.Sprintf("%s/%s", cdnConfig.SpacesConfig.SpacesUrl, key)
	oembedURL := fmt.Sprintf("%s/oembed/%s", cdnConfig.CdnEndpoint, key)
	res, err := http.Head(imageURL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fiber.NewError(res.StatusCode, "An error occurred redirecting to the image")
	}

	if queries.Download == "true" {
		ctx.Set("Content-Type", res.Header.Values("content-type")[0])
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, key))
		ctx.Set("Content-Length", fmt.Sprintf("%v", res.ContentLength))

		resp, err := http.Get(imageURL)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return ctx.Send(body)
	}

	if ctx.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
		ctx.Type("html")
		return ctx.Send([]byte(fmt.Sprintf(
			`<!DOCTYPE html>
			<html>
				<head>
					<meta name="theme-color" content="#dd9323">
					<meta property="og:title" content="%v">
					<meta content="%v" property="og:image">
					<link type="application/json+oembed" href="%v" />
				</head>
			</html>`,
			key, imageURL, oembedURL)),
		)
	} else {
		return ctx.Redirect(imageURL, fiber.StatusMovedPermanently)
	}
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

	respErr := DeleteFile(s, id)
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
	body := new(FolderPostRequest)

	if err := ctx.BodyParser(body); err != nil {
		respErr := NewResponseByError(fiber.StatusBadRequest, err)
		return ctx.JSON(respErr)
	}

	if body.Name == "" {
		respErr := NewResponse(fiber.StatusBadRequest, "Folder name required.")
		return ctx.JSON(respErr)
	}

	folder, respErr := NewFolder(body.Name)
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
