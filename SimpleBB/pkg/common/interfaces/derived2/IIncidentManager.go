package derived2

import (
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"net"
)

type IIncidentManager interface {
	Start() (err error)
	Stop() (err error)
	ReportIncident(itype cmi.IEnumValue, email simple.Email, userIPA net.IP)
}
