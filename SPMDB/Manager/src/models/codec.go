package models

type Codec struct {
	Id          ID     `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
