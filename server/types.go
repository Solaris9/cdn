package main

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

type ImageResponse struct {
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
