package models

type Actor struct {
	Id          ID     `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
