package main

import "time"

type User struct {
	UID   string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
	Admin bool   `json:"admin"`
}

type Config struct {
	SpacesConfig SpacesConfig
	CdnEndpoint  string
	AccessToken  string
	Production   bool
}

type SpacesConfig struct {
	SpacesAccessKey string
	SpacesSecretKey string
	SpacesEndpoint  string
	SpacesUrl       string
	SpacesCdn       string
	SpacesName      string
	SpacesRegion    string
}

type FileResult struct {
	CdnUrl       string    `json:"cdn_url"`
	SpacesUrl    string    `json:"spaces_url"`
	SpacesCdn    string    `json:"spaces_cdn"`
	FileName     string    `json:"file_name"`
	LastModified time.Time `json:"last_modified"`
	Size         int64     `json:"size"`
}

type FolderResult struct {
	CreateTime time.Time     `json:"create_time"`
	UpdateTime time.Time     `json:"update_time"`
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Files      []*FileResult `json:"files,omitempty"`
}

type FoldersResult struct {
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Size       int       `json:"size"`
}

type ImageResult struct {
	Url     string `json:"url"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
}

type DeletedImageResponse struct {
	Id      string `json:"id"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
}

type FolderPostRequest struct {
	Name string `json:"name"`
}

type FolderPatchRequest struct {
	Name   string   `json:"name"`
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

type Embed struct {
	Type         string `json:"type"`
	AuthorName   string `json:"author_name"`
	ProviderName string `json:"provider_name"`
}
