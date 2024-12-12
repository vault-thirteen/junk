package up

import (
	"database/sql"
	"errors"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UserRoles"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"time"
)

type userParameters struct {
	Id                cmb.Id        `json:"id"`
	PreRegTime        *time.Time    `json:"preRegTime,omitempty"`
	Email             simple.Email  `json:"email,omitempty"`
	Name              simple.Name   `json:"name,omitempty"`
	ApprovalTime      *time.Time    `json:"approvalTime,omitempty"`
	RegTime           *time.Time    `json:"regTime,omitempty"`
	Roles             *ur.UserRoles `json:"roles,omitempty"`
	LastBadLogInTime  *time.Time    `json:"lastBadLogInTime,omitempty"`
	BanTime           *time.Time    `json:"banTime,omitempty"`
	LastBadActionTime *time.Time    `json:"lastBadActionTime,omitempty"`
}

func NewUserParameters() (up base2.IUserParameters) {
	return &userParameters{
		Roles: &ur.UserRoles{},
	}
}

func NewUserParametersFromScannableSource(src cmi.IScannable) (up base2.IUserParameters, err error) {
	up = NewUserParameters()

	roles := up.GetRoles()

	err = src.Scan(
		up.GetIdPtr(),
		up.GetPreRegTimePtr(),
		up.GetEmailPtr(),
		up.GetNamePtr(),
		up.GetApprovalTimePtr(),
		up.GetRegTimePtr(),
		&roles.IsAuthor,
		&roles.IsWriter,
		&roles.IsReader,
		&roles.CanLogIn,
		up.GetLastBadLogInTimePtr(),
		up.GetBanTimePtr(),
		up.GetLastBadActionTimePtr(),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return up, nil
}

// Emulated class members.
func (up *userParameters) GetIdPtr() *cmb.Id                    { return &up.Id }
func (up *userParameters) GetPreRegTimePtr() **time.Time        { return &up.PreRegTime }
func (up *userParameters) GetEmailPtr() *simple.Email           { return &up.Email }
func (up *userParameters) GetNamePtr() *simple.Name             { return &up.Name }
func (up *userParameters) GetApprovalTimePtr() **time.Time      { return &up.ApprovalTime }
func (up *userParameters) GetRegTimePtr() **time.Time           { return &up.RegTime }
func (up *userParameters) GetRolesPtr() **ur.UserRoles          { return &up.Roles }
func (up *userParameters) GetLastBadLogInTimePtr() **time.Time  { return &up.LastBadLogInTime }
func (up *userParameters) GetBanTimePtr() **time.Time           { return &up.BanTime }
func (up *userParameters) GetLastBadActionTimePtr() **time.Time { return &up.LastBadActionTime }

func (up *userParameters) GetId() cmb.Id                    { return up.Id }
func (up *userParameters) GetPreRegTime() *time.Time        { return up.PreRegTime }
func (up *userParameters) GetEmail() simple.Email           { return up.Email }
func (up *userParameters) GetName() simple.Name             { return up.Name }
func (up *userParameters) GetApprovalTime() *time.Time      { return up.ApprovalTime }
func (up *userParameters) GetRegTime() *time.Time           { return up.RegTime }
func (up *userParameters) GetRoles() *ur.UserRoles          { return up.Roles }
func (up *userParameters) GetLastBadLogInTime() *time.Time  { return up.LastBadLogInTime }
func (up *userParameters) GetBanTime() *time.Time           { return up.BanTime }
func (up *userParameters) GetLastBadActionTime() *time.Time { return up.LastBadActionTime }

func (up *userParameters) SetId(id cmb.Id)              { up.Id = id }
func (up *userParameters) SetName(name simple.Name)     { up.Name = name }
func (up *userParameters) SetRoles(roles *ur.UserRoles) { up.Roles = roles }
