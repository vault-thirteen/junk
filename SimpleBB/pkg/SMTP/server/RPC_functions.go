package server

import (
	"fmt"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/rpc"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// RPC functions.

func (srv *Server) sendMessage(p *sm.SendMessageParams) (result *sm.SendMessageResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Recipient) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RecipientIsNotSet, RpcErrorMsg_RecipientIsNotSet, nil)
	}
	if len(p.Subject) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SubjectIsNotSet, RpcErrorMsg_SubjectIsNotSet, nil)
	}
	if len(p.Message) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIsNotSet, RpcErrorMsg_MessageIsNotSet, nil)
	}

	srv.mailerGuard.Lock()
	defer srv.mailerGuard.Unlock()

	err := srv.mailer.SendMail([]string{p.Recipient.ToString()}, p.Subject.ToString(), p.Message.ToString())
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MailerError, fmt.Sprintf(RpcErrorMsgF_MailerError, err.Error()), err)
	}

	result = &sm.SendMessageResult{}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *sm.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	trc, src := srv.js.GetRequestsCount()

	result = &sm.ShowDiagnosticDataResult{
		RequestsCount: cmr.RequestsCount{
			TotalRequestsCount:      cmb.Text(trc),
			SuccessfulRequestsCount: cmb.Text(src),
		},
	}

	return result, nil
}
