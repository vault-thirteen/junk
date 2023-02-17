package repository

// Storage описывает контракт взаимодействия
// с хранилищем микросервиса.
type Storage interface {
	BaseActions
	EventActions
}
