package let

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	enum "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Enum"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
)

type logEventType base.IEnum

const (
	LogEventType_LogIn   = 1
	LogEventType_LogOut  = 2 // Self logging out.
	LogEventType_LogOutA = 3 // Logging out by an administrator.

	LogEventTypeMax = LogEventType_LogOutA
)

func NewLogEventType() derived1.ILogEventType {
	return enum.NewEnumFast(ev.NewEnumValue(LogEventTypeMax))
}

func NewLogEventTypeWithValue(value base.IEnumValue) derived1.ILogEventType {
	let := NewLogEventType()
	let.SetValueFast(value)
	return let
}
