package main

type User struct {
	UID   string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
	Admin bool   `json:"admin"`
}

type File struct {
	Name  string `json:"name"`
	Ext   string `json:"ext"`
	Owner string `json:"owner"`
	Size  int64  `json:"size"`
}

type Folder struct {
	ID    string   `json:"id"`
	Owner string   `json:"owner"`
	Name  string   `json:"name"`
	Files []string `json:"-"`
}

type Config struct {
	SpacesConfig SpacesConfig
	CdnEndpoint  string
	AccessToken  string
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

type UploadedFile struct {
	ID   string
	Ext  string
	Size int64
}

type ImageResponse struct {
	Url     string `json:"url"`
	Success bool   `json:"success"`
}
