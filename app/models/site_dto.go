package models

type SingleApiConfig struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type MultiApiConfig struct {
	Urls []SingleApiConfig `json:"urls"`
}
