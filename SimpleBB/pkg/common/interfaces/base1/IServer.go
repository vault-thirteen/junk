package base

import (
	"sync"
	"sync/atomic"
)

type IServer interface {
	Start() error
	ReportStart()
	Stop() error
	GetStopChannel() *chan bool
	UseConstructor(ISettings) (IServer, error)
	GetSubRoutinesWG() *sync.WaitGroup
	GetMustStopAB() *atomic.Bool
}
