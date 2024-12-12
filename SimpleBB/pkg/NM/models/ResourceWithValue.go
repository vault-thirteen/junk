package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type ResourceWithValue struct {
	Id    cmb.Id `json:"id"`
	Value any    `json:"value"`
}
