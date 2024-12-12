package models

type Server struct {
	Id          ID     `json:"Id"`
	Address     string `json:"Address"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
