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
	Files   []*File
	Data    *FolderData
	Created time.Time
	Updates []firestore.Update
}

type FolderData struct {
	ID    string   `json:"id"`
	Owner string   `json:"owner"`
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

// creates a new folder
func NewFolder(name, owner string) (*Folder, *JSONResponse) {
	firebaseCtx := context.Background()
	folderID := randSeq(8)
	folderData := &FolderData{
		ID:    folderID,
		Owner: owner,
		Name:  name,
		Files: make([]string, 0),
	}

	doc, err := cdnFirestore.Collection("folders").Doc(folderID).Create(firebaseCtx, folderData)
	if err != nil {
		return nil, NewResponseByError(fiber.StatusInternalServerError, err)
	}

	folder := new(Folder)
	folder.Created = doc.UpdateTime
	folder.Data = folderData

	return folder, nil
}

// gets a folder, optionally cache all files in it
func FolderFor(id string, withFiles bool) (*Folder, *JSONResponse) {
	firebaseCtx := context.Background()
	folder := new(Folder)

	doc, err := cdnFirestore.Collection("folders").Doc(id).Get(firebaseCtx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, NewResponse(fiber.StatusNotFound, "Folder not found")
		}

		return nil, NewResponseByError(fiber.StatusInternalServerError, err)
	}

	folder.Created = doc.CreateTime

	folder.Data = new(FolderData)
	doc.DataTo(folder.Data)

	if withFiles {
		folder.GetFiles()
	}

	return folder, nil
}

// caches the files
func (folder *Folder) GetFiles() error {
	firebaseCtx := context.Background()

	docs, err := cdnFirestore.Collection("files").Where("ID", "in", folder.Data.Files).Documents(firebaseCtx).GetAll()

	if err != nil {
		return err
	}

	folder.Files = make([]*File, len(docs))

	for _, doc := range docs {
		fileData := new(FileData)
		doc.DataTo(fileData)

		file := &File{
			Data:    fileData,
			Created: doc.CreateTime,
		}

		folder.Files = append(folder.Files, file)
	}

	return nil
}

// returns a file
func (folder *Folder) GetFile(id string) *File {
	for _, file := range folder.Files {
		if file.Data.ID == id {
			return file
		}
	}

	firebaseCtx := context.Background()

	doc, err := cdnFirestore.Collection("files").Doc(id).Get(firebaseCtx)
	if err != nil && status.Code(err) == codes.NotFound {
		return nil
	} else if err == nil {
		fileData := new(FileData)
		doc.DataTo(fileData)

		file := &File{
			Data:    fileData,
			Created: doc.CreateTime,
		}

		folder.Data.Files = append(folder.Data.Files, id)
		folder.Files = append(folder.Files, file)

		return file
	}

	return nil
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

	if cacheFiles {
		folder.GetFiles()
	}
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

func (folder *Folder) CheckOwner(owner string) *JSONResponse {
	if folder.Data.Owner != owner {
		return NewResponse(fiber.StatusForbidden, "Cannot delete folder not owned.")
	}

	return nil
}

func (folder *Folder) Delete() *JSONResponse {
	firebaseCtx := context.Background()
	_, err := cdnFirestore.Collection("files").Doc(folder.Data.ID).Delete(firebaseCtx)

	if err != nil {
		return NewResponseByError(fiber.StatusInternalServerError, err)
	}

	return nil
}

// converts the folder object to a json object
func (folder *Folder) ToJSON() []byte {
	dataMap := NewJSONModifer(folder.Data)
	dataMap.AddField("date", folder.Created)
	return dataMap.ToJSON()
}

func (folder *Folder) Save() *JSONResponse {
	firebaseCtx := context.Background()
	_, err := cdnFirestore.Collection("folders").Doc(folder.Data.ID).Update(firebaseCtx, folder.Updates)

	if err != nil {
		return NewResponseByError(fiber.StatusInternalServerError, err)
	}

	return nil
}

func removeFile(files []*File, id string) []*File {
	var index = -1

	for i, file := range files {
		if file.Data.ID == id {
			index = i
		}
	}

	if index > 0 {
		return append(files[:index], files[index+1:]...)
	}

	return files
}
