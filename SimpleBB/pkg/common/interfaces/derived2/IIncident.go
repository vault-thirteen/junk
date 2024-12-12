package derived2

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"net"
	"time"
)

type IIncident interface {
	Check() (err error)

	// Emulated class members.
	GetId() (id cmb.Id)
	GetModule() (module derived1.IModule)
	GetType() (t derived1.IIncidentType)
	GetTime() (time time.Time)
	GetEmail() (email simple.Email)
	GetUserIPA() (userIPA net.IP)
}
