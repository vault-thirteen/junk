package rpc

import (
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

// Ping.

type PingParams = rpc2.PingParams
type PingResult = rpc2.PingResult

// IP address list.

type BlockIPAddressParams struct {
	rpc2.CommonParams

	// IP address of a client to block.
	UserIPA cm.IPAS `json:"userIPA"`

	// Block time in seconds.
	// This is a period for which the specified IP address will be blocked. If
	// for some reason a record with the specified IP address already exists,
	// this time will be added to an already existing value.
	BlockTimeSec base2.Count `json:"blockTimeSec"`
}
type BlockIPAddressResult = rpc2.CommonResultWithSuccess

type IsIPAddressBlockedParams struct {
	rpc2.CommonParams

	// IP address of a client to check.
	UserIPA cm.IPAS `json:"userIPA"`
}
type IsIPAddressBlockedResult struct {
	rpc2.CommonResult

	IsBlocked base2.Flag `json:"isBlocked"`
}

// Other.

type ShowDiagnosticDataParams struct{}
type ShowDiagnosticDataResult struct {
	rpc2.CommonResult
	rpc2.RequestsCount
}
