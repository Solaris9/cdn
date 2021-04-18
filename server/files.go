package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
)

type File struct {
	ID    string `json:"id"`
	Ext   string `json:"ext"`
	Owner string `json:"owner"`
	Size  int64  `json:"size"`
}

func UploadFile(s *session.Session, ctx *fiber.Ctx) (string, *JSONResponse) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return "", NewResponse(fiber.StatusInternalServerError, "Failed to get uploaded file.")
	}

	size := fileHeader.Size
	buffer := make([]byte, size)

	uploadedFile, err := fileHeader.Open()
	if err != nil {
		return "", NewResponse(fiber.StatusInternalServerError, "Failed to open uploaded file.")
	}

	_, err = uploadedFile.Read(buffer)
	if err != nil {
		return "", NewResponse(fiber.StatusInternalServerError, "Failed to read uploaded file.")
	}

	uploadedFile.Close()

	fileName := randSeq(8) + filepath.Ext(fileHeader.Filename)

	object := s3.PutObjectInput{
		Bucket:               aws.String(cdnConfig.SpacesConfig.SpacesName),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("public-read"),
		Body:                 strings.NewReader(string(buffer)),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ServerSideEncryption: aws.String("AES256"),
	}

	_, err = s3.New(s).PutObject(&object)
	if err != nil {
		return "", NewResponseByError(fiber.StatusInternalServerError, err)
	}

	return fileName, nil
}

// func CheckOwner(s *session.Session, file, owner string) *JSONResponse {
// 	object := s3.HeadObjectInput{
// 		Bucket: aws.String(cdnConfig.SpacesConfig.SpacesName),
// 		Key:    aws.String(file),
// 	}

// 	out, err := s3.New(s).HeadObject(&object)
// 	if err != nil {
// 		return NewResponseByError(fiber.StatusInternalServerError, err)
// 	}

// 	if *out.Metadata["Owner"] != owner {
// 		return NewResponse(fiber.StatusForbidden, "Cannot delete file not owned.")
// 	}

// 	return nil
// }

func DeleteFile(s *session.Session, file string) *JSONResponse {
	object := &s3.DeleteObjectInput{
		Bucket: aws.String(cdnConfig.SpacesConfig.SpacesName),
		Key:    aws.String(file),
	}

	_, err := s3.New(s).DeleteObject(object)
	if err != nil {
		return NewResponseByError(fiber.StatusInternalServerError, err)
	}

	return nil
}

func GetFiles(s *session.Session) ([]*FileResult, error) {
	var files []*FileResult
	var shouldContinue = true
	var nextToken = ""

	for shouldContinue {
		objects, err := s3.New(s).ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:            aws.String(cdnConfig.SpacesConfig.SpacesName),
			ContinuationToken: aws.String(nextToken),
		})
		if err != nil {
			return nil, err
		}

		var data []*FileResult
		for _, obj := range objects.Contents {
			data = append(data, &FileResult{
				CdnUrl:       fmt.Sprintf("%v/%v", cdnConfig.CdnEndpoint, *obj.Key),
				SpacesUrl:    fmt.Sprintf("%v/%v", cdnConfig.SpacesConfig.SpacesUrl, *obj.Key),
				SpacesCdn:    fmt.Sprintf("%v/%v", cdnConfig.SpacesConfig.SpacesCdn, *obj.Key),
				FileName:     *obj.Key,
				LastModified: *obj.LastModified,
				Size:         *obj.Size,
			})
		}

		files = append(files, data...)
		if !*objects.IsTruncated {
			shouldContinue = false
			nextToken = ""
		} else {
			nextToken = *objects.NextContinuationToken
		}
	}

	return files, nil
}

func GetFilesByKeys(s *session.Session, keys []string) ([]*FileResult, error) {
	var files []*FileResult
	var client = s3.New(s)

	for _, key := range keys {
		obj, err := client.HeadObject(&s3.HeadObjectInput{
			Bucket: &cdnConfig.SpacesConfig.SpacesName,
			Key:    &key,
		})

		if err != nil {
			return nil, err
		}

		files = append(files, &FileResult{
			CdnUrl:       fmt.Sprintf("%v/%v", cdnConfig.CdnEndpoint, key),
			SpacesUrl:    fmt.Sprintf("%v/%v", cdnConfig.SpacesConfig.SpacesUrl, key),
			SpacesCdn:    fmt.Sprintf("%v/%v", cdnConfig.SpacesConfig.SpacesCdn, key),
			FileName:     key,
			LastModified: *obj.LastModified,
			Size:         *obj.ContentLength,
		})
	}

	return files, nil
}
