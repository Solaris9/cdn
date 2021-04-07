package main

import (
	"context"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go"
)

func authorize(ctx *fiber.Ctx) error {
	authorization := ctx.Get("Authorization")
	if authorization == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "No authorization token provided.")
	}

	_ctx := context.Background()

	token, err := cdnAuth.VerifyIDToken(_ctx, authorization)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization token provided.")
	}

	ctx.Locals("user", token.UID)

	return ctx.Next()
}

func uploadFile(owner string, fileHeader *multipart.FileHeader) (string, error) {
	size := fileHeader.Size
	buffer := make([]byte, size)

	file, headerOpenErr := fileHeader.Open()
	if headerOpenErr != nil {
		return "", headerOpenErr
	}

	_, fileReadErr := file.Read(buffer)
	if fileReadErr != nil {
		return "", fileReadErr
	}

	file.Close()

	contentType := http.DetectContentType(buffer)
	fileID := randSeq(7)
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fileID + fileExt

	_, uploadErr := cdnSpaces.PutObject(
		cdnConfig.SpacesConfig.SpacesName,
		fileName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType:  contentType,
			UserMetadata: map[string]string{"x-amz-acl": "public-read"},
		},
	)

	if uploadErr != nil {
		return "", uploadErr
	}

	createFileRecord(fileID, fileExt, owner, size)

	return fileID, nil
}

func createFileRecord(ID string, ext string, owner string, size int64) {
	firebaseCtx := context.Background()

	cdnFirestore.Collection("files").Doc(ID).Set(firebaseCtx, File{
		Name:  ID,
		Ext:   ext,
		Owner: owner,
		Size:  size,
	})
}

func getFileRecord(id string) (*File, error) {
	firebaseCtx := context.Background()

	doc, err := cdnFirestore.Collection("files").Doc(id).Get(firebaseCtx)
	if err != nil {
		return nil, err
	}

	file := new(File)
	doc.DataTo(file)

	return file, nil
}

func getFolderRecord(id string) (*Folder, []*File, error) {
	firebaseCtx := context.Background()

	doc, err := cdnFirestore.Collection("folders").Doc("discord").Get(firebaseCtx)
	if err != nil {
		return nil, nil, err
	}

	var files []*File
	folder := new(Folder)
	doc.DataTo(folder)

	docs, err := cdnFirestore.Collection("files").Where("id", "in", folder.Files).Documents(firebaseCtx).GetAll()

	if err != nil {
		return nil, nil, err
	}

	for _, doc := range docs {
		file := new(File)
		doc.DataTo(file)
		files = append(files, file)
	}

	return folder, files, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
