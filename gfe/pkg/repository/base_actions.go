package repository

import "io"

// BaseActions описывает контракт для базовых действий
// с репозиториями.
type BaseActions interface {
	io.Closer

	// Open начинает работу со хранилищем.
	Open()

	// Wait дожидается момента когда хранилище станет доступно.
	Wait()

	// Ping проверяет связь с сервером базы данных.
	Ping() error

	// IsReady читает состояние готовности хранилища.
	IsReady() bool
}
