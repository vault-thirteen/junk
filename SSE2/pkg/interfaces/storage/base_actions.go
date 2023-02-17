package storage

type BaseActions interface {
	Ping() error
	IsReady() bool
	Wait()
}
