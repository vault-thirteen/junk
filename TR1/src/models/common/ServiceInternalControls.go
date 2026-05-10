package cm

import (
	"sync"
	"sync/atomic"

	"github.com/vault-thirteen/TR1/src/libraries/avm"
)

type ServiceInternalControls struct {
	subRoutines *sync.WaitGroup
	ssp         *avm.SSP

	// A flag showing that service must stop itself.
	mustStop *atomic.Bool
}

func NewInternalControls() *ServiceInternalControls {
	return &ServiceInternalControls{
		subRoutines: &sync.WaitGroup{},
		ssp:         avm.NewSSP(),
		mustStop:    &atomic.Bool{},
	}
}
