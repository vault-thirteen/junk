package models

type Tag struct {
	Id          ID     `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Colour      string `json:"Colour"`
}
