package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileData struct {
	ID    string `json:"id"`
	Ext   string `json:"ext"`
	Owner string `json:"owner"`
	Size  int64  `json:"size"`
}

type File struct {
	Data    *FileData
	Created time.Time
}

func (file *File) URL() string {
	return fmt.Sprintf("%v/%v", cdnConfig.CdnEndpoint, file.Data.ID)
}

func (file *File) ImageURL() string {
	return fmt.Sprintf("%v/%v%v", cdnConfig.SpacesConfig.SpacesUrl, file.Data.ID, file.Data.Ext)
}

func (file *File) Name() string {
	return fmt.Sprintf("%v%v", file.Data.ID, file.Data.Ext)
}

func (file *File) ToJSON() []byte {
	dataMap := NewJSONModifer(file.Data)
	dataMap.AddField("date", file.Created)
	dataMap.AddField("url", file.ImageURL())
	return dataMap.ToJSON()
}

func FileFor(id string) (*File, *JSONResponse) {
	firebaseCtx := context.Background()

	doc, err := cdnFirestore.Collection("files").Doc(id).Get(firebaseCtx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, NewResponse(fiber.StatusNotFound, "File not found")
		}

		return nil, NewResponseByError(fiber.StatusInternalServerError, err)
	}

	data := new(FileData)
	doc.DataTo(data)

	file := &File{
		Data:    data,
		Created: doc.CreateTime,
	}

	return file, nil
}

func UploadFile(ctx *fiber.Ctx) (*File, *JSONResponse) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, "Failed to get uploaded file.")
	}

	size := fileHeader.Size
	buffer := make([]byte, size)

	uploadedFile, err := fileHeader.Open()
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, "Failed to open uploaded file.")
	}

	_, err = uploadedFile.Read(buffer)
	if err != nil {
		return nil, NewResponse(fiber.StatusInternalServerError, "Failed to read uploaded file.")
	}

	uploadedFile.Close()

	contentType := http.DetectContentType(buffer)

	fileID := randSeq(8)
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fileID + fileExt
	fileOwner := ctx.Locals("user").(string)

	options := minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	}

	_, err = cdnSpaces.PutObject(
		cdnConfig.SpacesConfig.SpacesName,
		fileName,
		uploadedFile,
		size,
		options,
	)

	if err != nil {
		return nil, NewResponseByError(fiber.StatusInternalServerError, err)
	}

	file, respErr := createFileRecord(fileID, fileExt, fileOwner, size)
	if err != nil {
		return nil, respErr
	}

	return file, nil
}

func (file *File) CheckOwner(owner string) *JSONResponse {
	if file.Data.Owner != owner {
		return NewResponse(fiber.StatusForbidden, "Cannot delete file not owned.")
	}

	return nil
}

func (file *File) Delete() *JSONResponse {
	err := cdnSpaces.RemoveObject(cdnConfig.SpacesConfig.SpacesName, file.Name())
	if err != nil {
		return NewResponseByError(fiber.StatusInternalServerError, err)
	}

	respErr := deleteFileRecord(file.Data.ID)
	if err != nil {
		return respErr
	}

	return nil
}

func createFileRecord(fileID, fileExt, fileOwner string, fileSize int64) (*File, *JSONResponse) {
	firebaseCtx := context.Background()
	data := &FileData{
		ID:    fileID,
		Ext:   fileExt,
		Owner: fileOwner,
		Size:  fileSize,
	}

	res, err := cdnFirestore.Collection("files").Doc(fileID).Set(firebaseCtx, data)

	if err != nil {
		return nil, NewResponseByError(fiber.StatusInternalServerError, err)
	}

	file := &File{
		Data:    data,
		Created: res.UpdateTime,
	}

	return file, nil
}

func deleteFileRecord(id string) *JSONResponse {
	firebaseCtx := context.Background()
	_, err := cdnFirestore.Collection("files").Doc(id).Delete(firebaseCtx)

	if err != nil {
		return NewResponseByError(fiber.StatusInternalServerError, err)
	}

	return nil
}

func deleteFileRecords(ids ...string) *JSONResponse {
	firebaseCtx := context.Background()

	batch := cdnFirestore.Batch()
	docs, err := cdnFirestore.Collection("files").Where("ID", "in", ids).Documents(firebaseCtx).GetAll()

	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}

	_, err = batch.Commit(firebaseCtx)
	if err != nil {
		return NewResponseByError(fiber.StatusInternalServerError, err)
	}

	return nil
}
