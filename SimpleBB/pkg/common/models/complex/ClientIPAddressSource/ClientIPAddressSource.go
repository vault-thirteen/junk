package cipas

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	enum "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Enum"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
)

type clientIPAddressSource base.IEnum

const (
	ClientIPAddressSource_Direct       = 1
	ClientIPAddressSource_CustomHeader = 2

	ClientIPAddressSourceMax = ClientIPAddressSource_CustomHeader
)

func NewClientIPAddressSource() derived1.IClientIPAddressSource {
	return enum.NewEnumFast(ev.NewEnumValue(ClientIPAddressSourceMax))
}

func NewClientIPAddressSourceWithValue(value base.IEnumValue) derived1.IClientIPAddressSource {
	cipas := NewClientIPAddressSource()
	cipas.SetValueFast(value)
	return cipas
}
