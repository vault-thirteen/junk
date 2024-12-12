package rpc

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type RequestsCount struct {
	TotalRequestsCount      cmb.Text `json:"totalRequestsCount"`
	SuccessfulRequestsCount cmb.Text `json:"successfulRequestsCount"`
}
