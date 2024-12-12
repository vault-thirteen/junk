package inc

import (
	"errors"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/IncidentType"
	m "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Module"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"net"
	"time"
)

const (
	ErrIncidentIsNotSet = "incident is not set"
)

type incident struct {
	Id      cmb.Id
	Module  derived1.IModule
	Type    derived1.IIncidentType
	Time    time.Time
	Email   simple.Email
	UserIPA net.IP
}

func NewIncident() derived2.IIncident {
	return &incident{
		Module: m.NewModule(),
		Type:   it.NewIncidentType(),
	}
}

func NewIncidentWithFields(itype cmi.IEnumValue, email simple.Email, userIPA net.IP) derived2.IIncident {
	return &incident{
		Module:  m.NewModule(),
		Type:    it.NewIncidentTypeWithValue(itype),
		Time:    time.Now(),
		Email:   email,
		UserIPA: userIPA,
	}
}

func (i *incident) Check() (err error) {
	if i == nil {
		return errors.New(ErrIncidentIsNotSet)
	}

	return nil
}

// Emulated class members.
func (i *incident) GetId() (id cmb.Id) {
	return i.Id
}
func (i *incident) GetModule() (module derived1.IModule) { return i.Module }
func (i *incident) GetType() (t derived1.IIncidentType)  { return i.Type }
func (i *incident) GetTime() (time time.Time)            { return i.Time }
func (i *incident) GetEmail() (email simple.Email)       { return i.Email }
func (i *incident) GetUserIPA() (userIPA net.IP)         { return i.UserIPA }
