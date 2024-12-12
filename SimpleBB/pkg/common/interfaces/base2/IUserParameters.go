package base2

import (
	ur "github.com/vault-thirteen/SimpleBB/pkg/common/models/UserRoles"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"time"
)

type IUserParameters interface {
	// Emulated class members.
	GetIdPtr() *cmb.Id
	GetPreRegTimePtr() **time.Time
	GetEmailPtr() *simple.Email
	GetNamePtr() *simple.Name
	GetApprovalTimePtr() **time.Time
	GetRegTimePtr() **time.Time
	GetRolesPtr() **ur.UserRoles
	GetLastBadLogInTimePtr() **time.Time
	GetBanTimePtr() **time.Time
	GetLastBadActionTimePtr() **time.Time

	GetId() cmb.Id
	GetPreRegTime() *time.Time
	GetEmail() simple.Email
	GetName() simple.Name
	GetApprovalTime() *time.Time
	GetRegTime() *time.Time
	GetRoles() *ur.UserRoles
	GetLastBadLogInTime() *time.Time
	GetBanTime() *time.Time
	GetLastBadActionTime() *time.Time

	SetId(id cmb.Id)
	SetName(name simple.Name)
	SetRoles(roles *ur.UserRoles)
}
