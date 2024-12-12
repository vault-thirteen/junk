package json

import (
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
)

func ToJson(x cmi.IToString) []byte {
	return []byte(x.ToString())
}
