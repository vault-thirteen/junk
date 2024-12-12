package set

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	enum "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Enum"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
)

type systemEventType base.IEnum

const (
	SystemEventType_ThreadParentChange    = 1 // -> Users subscribed to the thread.
	SystemEventType_ThreadNameChange      = 2 // -> Users subscribed to the thread.
	SystemEventType_ThreadDeletion        = 3 // -> Users subscribed to the thread.
	SystemEventType_ThreadNewMessage      = 4 // -> Users subscribed to the thread.
	SystemEventType_ThreadMessageEdit     = 5 // -> Users subscribed to the thread.
	SystemEventType_ThreadMessageDeletion = 6 // -> Users subscribed to the thread.
	SystemEventType_MessageTextEdit       = 7 // -> Author of the message.
	SystemEventType_MessageParentChange   = 8 // -> Author of the message.
	SystemEventType_MessageDeletion       = 9 // -> Author of the message.

	SystemEventTypeMax = SystemEventType_MessageDeletion
)

func NewSystemEventType() derived1.ISystemEventType {
	return enum.NewEnumFast(ev.NewEnumValue(SystemEventTypeMax))
}

func NewSystemEventTypeWithValue(value base.IEnumValue) derived1.ISystemEventType {
	let := NewSystemEventType()
	let.SetValueFast(value)
	return let
}
