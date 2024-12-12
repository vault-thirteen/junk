package server

// RPC handlers.

import (
	"encoding/json"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/rpc"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/models/settings"
)

func (srv *Server) initRpc() (err error) {
	rpcDurationFieldName := cs.RpcDurationFieldName
	rpcRequestIdFieldName := cs.RpcRequestIdFieldName

	ps := &jrm1.ProcessorSettings{
		CatchExceptions:    true,
		LogExceptions:      true,
		CountRequests:      true,
		DurationFieldName:  &rpcDurationFieldName,
		RequestIdFieldName: &rpcRequestIdFieldName,
	}

	srv.js, err = jrm1.NewProcessor(ps)
	if err != nil {
		return err
	}

	fns := []jrm1.RpcFunction{
		srv.Ping,
		srv.CreateCaptcha,
		srv.CheckCaptcha,
		srv.ShowDiagnosticData,
	}

	for _, fn := range fns {
		err = srv.js.AddFunc(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

// Ping.

func (srv *Server) Ping(_ *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	result = rm.PingResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// Captcha.

func (srv *Server) CreateCaptcha(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.CreateCaptchaParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.CreateCaptchaResult
	r, re = srv.createCaptcha()
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) CheckCaptcha(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.CheckCaptchaParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.CheckCaptchaResult
	r, re = srv.checkCaptcha(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Other.

func (srv *Server) ShowDiagnosticData(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ShowDiagnosticDataParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ShowDiagnosticDataResult
	r, re = srv.showDiagnosticData()
	if re != nil {
		return nil, re
	}

	return r, nil
}
