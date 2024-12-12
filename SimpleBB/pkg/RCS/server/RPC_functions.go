package server

import (
	"encoding/base64"
	"fmt"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	rc "github.com/vault-thirteen/RingCaptcha/server"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/rpc"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// RPC functions.

func (srv *Server) createCaptcha() (result *rm.CreateCaptchaResult, re *jrm1.RpcError) {
	ccResponse, err := srv.captchaManager.CreateCaptcha()
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CreateError, fmt.Sprintf(RpcErrorMsgF_CreateError, err.Error()), nil)
	}

	result = &rm.CreateCaptchaResult{
		TaskId:              ccResponse.TaskId,
		ImageFormat:         ccResponse.ImageFormat,
		IsImageDataReturned: ccResponse.IsImageDataReturned,
	}

	if ccResponse.IsImageDataReturned {
		result.ImageDataB64 = base64.StdEncoding.EncodeToString(ccResponse.ImageData)
	}

	return result, nil
}

func (srv *Server) checkCaptcha(p *rm.CheckCaptchaParams) (result *rm.CheckCaptchaResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.TaskId) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TaskIdIsNotSet, RpcErrorMsg_TaskIdIsNotSet, nil)
	}
	if p.Value == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_AnswerIsNotSet, RpcErrorMsg_AnswerIsNotSet, nil)
	}

	resp, err := srv.captchaManager.CheckCaptcha(&rc.CheckCaptchaRequest{TaskId: p.TaskId, Value: p.Value})
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CheckError, fmt.Sprintf(RpcErrorMsgF_CheckError, err.Error()), nil)
	}

	result = &rm.CheckCaptchaResult{
		TaskId:    p.TaskId,
		IsSuccess: resp.IsSuccess,
	}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *rm.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	trc, src := srv.js.GetRequestsCount()

	result = &rm.ShowDiagnosticDataResult{
		RequestsCount: cmr.RequestsCount{
			TotalRequestsCount:      cmb.Text(trc),
			SuccessfulRequestsCount: cmb.Text(src),
		},
	}

	return result, nil
}
