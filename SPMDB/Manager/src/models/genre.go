package models

type Genre struct {
	Id          ID     `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
