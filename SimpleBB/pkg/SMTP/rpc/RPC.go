package rpc

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

// Ping.

type PingParams = rpc2.PingParams
type PingResult = rpc2.PingResult

// Message.

type SendMessageParams struct {
	Recipient cm.Email `json:"recipient"`
	Subject   cmb.Text `json:"subject"`
	Message   cmb.Text `json:"message"`
}
type SendMessageResult struct {
	rpc2.CommonResult
}

// Other.

type ShowDiagnosticDataParams struct{}
type ShowDiagnosticDataResult struct {
	rpc2.CommonResult
	rpc2.RequestsCount
}
