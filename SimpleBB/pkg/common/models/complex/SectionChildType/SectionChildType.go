package sct

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	enum "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Enum"
	ev "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
)

type sectionChildType base.IEnum

const (
	SectionChildType_Section = 1
	SectionChildType_Forum   = 2
	SectionChildType_None    = 3

	SectionChildTypeMax = SectionChildType_None
)

func NewSectionChildType() derived1.ISectionChildType {
	return enum.NewEnumFast(ev.NewEnumValue(SectionChildTypeMax))
}

func NewSectionChildTypeWithValue(value base.IEnumValue) derived1.ISectionChildType {
	sct := NewSectionChildType()
	sct.SetValueFast(value)
	return sct
}
