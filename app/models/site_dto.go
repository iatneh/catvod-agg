package models

type SingleApiConfig struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type MultiApiConfig struct {
	Urls []SingleApiConfig `json:"urls"`
}
