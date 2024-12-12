package rt

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Enum"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
)

type resourceType base.IEnum

const (
	ResourceType_Text   = 1
	ResourceType_Number = 2

	ResourceTypeMax = ResourceType_Number
)

func NewResourceType() derived1.IResourceType {
	return enum.NewEnumFast(ev.NewEnumValue(ResourceTypeMax))
}

func NewResourceTypeWithValue(value base.IEnumValue) derived1.IResourceType {
	rt := NewResourceType()
	rt.SetValueFast(value)
	return rt
}
