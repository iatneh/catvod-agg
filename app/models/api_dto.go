package models

type ApiJson struct {
	Url  string
	Name string
}

type ApiContent struct {
	Urls string
	Url  *ApiJson
}
