package derived2

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"time"
)

type ISystemEvent interface {
	// Emulated class members.
	GetIdPtr() (id *cmb.Id)
	GetTimePtr() (t *time.Time)
	GetSystemEventData() (sed ISystemEventData)
}
