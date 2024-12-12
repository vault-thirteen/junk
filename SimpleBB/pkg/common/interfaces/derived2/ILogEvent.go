package derived2

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"net"
	"time"
)

type ILogEvent interface {
	// Emulated class members.
	GetId() (id cmb.Id)
	GetTime() (time time.Time)
	GetType() (t derived1.ILogEventType)
	GetUserId() (userId cmb.Id)
	GetEmail() (email cms.Email)
	GetUserIPAB() (userIPAB net.IP)
	GetAdminId() (adminId *cmb.Id)
}
