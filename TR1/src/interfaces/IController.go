package interfaces

import (
	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/libraries/scheduler"
)

type IController interface {
	GetRpcFunctions() []jrm1.RpcFunction
	GetScheduledFunctions() []sch.ScheduledFn
	GetErrorsChan() (errorsChan *chan error)
	LinkWithService(service IService) (err error)
}
