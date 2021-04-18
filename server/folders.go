package main

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Folder struct {
	Files      []*File
	Data       *FolderData
	CreateTime time.Time
	UpdateTime time.Time
	Updates    []firestore.Update
}

type FolderData struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

// creates a new folder
func NewFolder(name string) (*Folder, *JSONResponse) {
	firebaseCtx := context.Background()
	folderID := randSeq(8)
	folderData := &FolderData{
		ID:    folderID,
		Name:  name,
		Files: make([]string, 0),
	}

	doc, err := cdnFirestore.Collection("folders").Doc(folderID).Create(firebaseCtx, folderData)
	if err != nil {
		return nil, NewResponseByError(fiber.StatusInternalServerError, err)
	}

	folder := new(Folder)
	folder.CreateTime = doc.UpdateTime
	folder.UpdateTime = doc.UpdateTime
	folder.Data = folderData

	return folder, nil
}

// gets a folder, optionally cache all files in it
func FolderFor(id string) (*Folder, *JSONResponse) {
	firebaseCtx := context.Background()
	folder := new(Folder)

	doc, err := cdnFirestore.Collection("folders").Doc(id).Get(firebaseCtx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, NewResponse(fiber.StatusNotFound, "Folder not found")
		}

		return nil, NewResponseByError(fiber.StatusInternalServerError, err)
	}

	folder.CreateTime = doc.UpdateTime
	folder.UpdateTime = doc.UpdateTime

	folder.Data = new(FolderData)
	doc.DataTo(folder.Data)

	return folder, nil
}

func (folder *Folder) IsChanged() bool {
	return len(folder.Updates) > 0
}

// a list of ids to add optionally cache all files again
func (folder *Folder) SetName(name string) {
	folder.Data.Name = name

	folder.Updates = append(folder.Updates, firestore.Update{
		Path:  "Name",
		Value: folder.Data.Name,
	})
}

// a list of ids to add optionally cache all files again
func (folder *Folder) AddFiles(files []string, cacheFiles bool) {
	folder.Data.Files = Set(append(folder.Data.Files, files...))

	folder.Updates = append(folder.Updates, firestore.Update{
		Path:  "Files",
		Value: folder.Data.Files,
	})
}

// a list of ids to remove
func (folder *Folder) RemoveFiles(files []string) {
	for _, file := range files {
		if index := indexOf(folder.Data.Files, file); index != -1 {
			folder.Data.Files = append(folder.Data.Files[:index], folder.Data.Files[index+1:]...)
		}
	}

	folder.Updates = append(folder.Updates, firestore.Update{
		Path:  "Files",
		Value: folder.Data.Files,
	})
}

// func (folder *Folder) CheckOwner(owner string) *JSONResponse {
// 	if folder.Data.Owner != owner {
// 		return NewResponse(fiber.StatusForbidden, "Cannot delete folder not owned.")
// 	}

// 	return nil
// }

func (folder *Folder) Delete() *JSONResponse {
	firebaseCtx := context.Background()
	_, err := cdnFirestore.Collection("folders").Doc(folder.Data.ID).Delete(firebaseCtx)

	if err != nil {
		return NewResponseByError(fiber.StatusInternalServerError, err)
	}

	return nil
}

// converts the folder object to a json object
func (folder *Folder) ToJSON() *JSONModifier {
	dataMap := NewJSONModifer(folder.Data)
	dataMap.AddField("created", folder.CreateTime)
	dataMap.AddField("updated", folder.UpdateTime)
	return dataMap
}

// converts the folder object to a json object
func (folder *Folder) ToJSONMap() map[string]interface{} {
	return folder.ToJSON().Map
}

func (folder *Folder) Save() *JSONResponse {
	firebaseCtx := context.Background()
	_, err := cdnFirestore.Collection("folders").Doc(folder.Data.ID).Update(firebaseCtx, folder.Updates)

	if err != nil {
		return NewResponseByError(fiber.StatusInternalServerError, err)
	}

	return nil
}

func removeFile(files []string, id string) []string {
	var index = -1

	for i, file := range files {
		if file == id {
			index = i
		}
	}

	if index > 0 {
		return append(files[:index], files[index+1:]...)
	}

	return files
}
