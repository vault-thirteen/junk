package server

import (
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	gm "github.com/vault-thirteen/SimpleBB/pkg/GWM/rpc"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/models/net"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

// RPC functions.

func (srv *Server) blockIPAddress(p *gm.BlockIPAddressParams) (result *gm.BlockIPAddressResult, re *jrm1.RpcError) {
	if !srv.settings.GetSystemSettings().GetIsFirewallUsed() {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_FirewallIsDisabled, RpcErrorMsg_FirewallIsDisabled, nil)
	}

	// Check parameters.
	if len(p.UserIPA) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IPAddressIsNotSet, RpcErrorMsg_IPAddressIsNotSet, nil)
	}

	if p.BlockTimeSec == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_BlockTimeIsNotSet, RpcErrorMsg_BlockTimeIsNotSet, nil)
	}

	userIPAB, err := cn.ParseIPA(p.UserIPA)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Search for an existing record.
	var n base2.Count
	n, err = srv.dbo.CountBlocksByIPAddress(userIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Insert a new record when none is found, or update an existing block.
	if n == 0 {
		err = srv.dbo.InsertBlock(userIPAB, p.BlockTimeSec)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	} else {
		err = srv.dbo.IncreaseBlockDuration(userIPAB, p.BlockTimeSec)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &gm.BlockIPAddressResult{
		Success: rpc2.Success{
			OK: true,
		},
	}
	return result, nil
}

func (srv *Server) isIPAddressBlocked(p *gm.IsIPAddressBlockedParams) (result *gm.IsIPAddressBlockedResult, re *jrm1.RpcError) {
	if !srv.settings.GetSystemSettings().GetIsFirewallUsed() {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_FirewallIsDisabled, RpcErrorMsg_FirewallIsDisabled, nil)
	}

	// Check parameters.
	if len(p.UserIPA) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IPAddressIsNotSet, RpcErrorMsg_IPAddressIsNotSet, nil)
	}

	userIPAB, err := cn.ParseIPA(p.UserIPA)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Search for an existing record.
	var n base2.Count
	n, err = srv.dbo.CountBlocksByIPAddress(userIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &gm.IsIPAddressBlockedResult{IsBlocked: n > 0}, nil
}

func (srv *Server) showDiagnosticData() (result *gm.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	trc, src := srv.js.GetRequestsCount()

	result = &gm.ShowDiagnosticDataResult{
		RequestsCount: rpc2.RequestsCount{
			TotalRequestsCount:      base2.Text(trc),
			SuccessfulRequestsCount: base2.Text(src),
		},
	}

	return result, nil
}
