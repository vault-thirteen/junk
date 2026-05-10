package interfaces

import (
	"sync"
	"sync/atomic"
)

type IService interface {
	Start() (err error)
	GetStopChannel() *chan bool
	ReportStart()
	Stop() (err error)
	GetSubRoutinesWG() *sync.WaitGroup
	GetMustStopAB() *atomic.Bool
	GetComponentByIndex(idx int) IServiceComponent
	GetConfiguration() IConfiguration
}
