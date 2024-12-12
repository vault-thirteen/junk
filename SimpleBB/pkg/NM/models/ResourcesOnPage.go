package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type ResourcesOnPage struct {
	ResourceIds []cmb.Id      `json:"resourceIds"`
	PageData    *cmr.PageData `json:"pageData,omitempty"`
}

func NewResourcesOnPage() (rop *ResourcesOnPage) {
	rop = &ResourcesOnPage{}
	return rop
}
