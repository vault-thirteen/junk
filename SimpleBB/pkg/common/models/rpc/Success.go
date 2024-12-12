package rpc

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type Success struct {
	OK cmb.Flag `json:"ok"`
}
