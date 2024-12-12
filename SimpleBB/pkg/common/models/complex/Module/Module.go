package m

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	enum "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Enum"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
)

type module base.IEnum

const (
	Module_ACM  = 1
	Module_GWM  = 2
	Module_MM   = 3
	Module_NM   = 4
	Module_RCS  = 5
	Module_SM   = 6
	Module_SMTP = 7

	ModuleMax = Module_SMTP
)

func NewModule() derived1.IModule {
	return enum.NewEnumFast(ev.NewEnumValue(ModuleMax))
}

func NewModuleWithValue(value base.IEnumValue) derived1.IModule {
	m := NewModule()
	m.SetValueFast(value)
	return m
}
