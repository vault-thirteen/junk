package it

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	enum "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Enum"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
)

type incidentType base.IEnum

const (
	IncidentType_IllegalAccessAttempt            = 1
	IncidentType_FakeToken                       = 2
	IncidentType_VerificationCodeMismatch        = 3
	IncidentType_DoubleLogInAttempt              = 4
	IncidentType_PreSessionHacking               = 5
	IncidentType_CaptchaAnswerMismatch           = 6
	IncidentType_PasswordMismatch                = 7
	IncidentType_PasswordChangeHacking           = 8
	IncidentType_EmailChangeHacking              = 9
	IncidentType_FakeIPA                         = 10
	IncidentType_ReadingNotificationOfOtherUsers = 11
	IncidentType_WrongDKey                       = 12

	IncidentTypeMax = IncidentType_WrongDKey
)

func NewIncidentType() derived1.IIncidentType {
	return enum.NewEnumFast(ev.NewEnumValue(IncidentTypeMax))
}

func NewIncidentTypeWithValue(value base.IEnumValue) derived1.IIncidentType {
	it := NewIncidentType()
	it.SetValueFast(value)
	return it
}
