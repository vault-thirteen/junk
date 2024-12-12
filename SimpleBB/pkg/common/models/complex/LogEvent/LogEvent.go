package le

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/LogEventType"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"net"
	"time"
)

type logEvent struct {
	Id       cmb.Id
	Time     time.Time
	Type     derived1.ILogEventType
	UserId   cmb.Id
	Email    cms.Email
	UserIPAB net.IP

	// ID of administrator for those events which were started by an
	// administrator.
	AdminId *cmb.Id
}

func NewLogEvent() derived2.ILogEvent {
	return &logEvent{
		Type: let.NewLogEventType(),
	}
}

func NewLogEventWithFields(t derived1.ILogEventType, userId cmb.Id, email cms.Email, userIPAB net.IP, adminId *cmb.Id) derived2.ILogEvent {
	return &logEvent{
		Type:     t,
		UserId:   userId,
		Email:    email,
		UserIPAB: userIPAB,
		AdminId:  adminId,
	}
}

// Emulated class members.
func (le *logEvent) GetId() (id cmb.Id) {
	return le.Id
}
func (le *logEvent) GetTime() (time time.Time) {
	return le.Time
}
func (le *logEvent) GetType() (t derived1.ILogEventType) {
	return le.Type
}
func (le *logEvent) GetUserId() (userId cmb.Id) {
	return le.UserId
}
func (le *logEvent) GetEmail() (email cms.Email) {
	return le.Email
}
func (le *logEvent) GetUserIPAB() (userIPAB net.IP) {
	return le.UserIPAB
}
func (le *logEvent) GetAdminId() (adminId *cmb.Id) {
	return le.AdminId
}
