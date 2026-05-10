package c

import (
	"encoding/json"
	"fmt"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) SendEmailMessage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.SendEmailMessageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.SendEmailMessageResult
	r, re = c.sendEmailMessage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) sendEmailMessage(p *rm.SendEmailMessageParams) (result *rm.SendEmailMessageResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Recipient) == 0 {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_EmailRecipientIsNotSet, rme.Msg_EmailRecipientIsNotSet, nil)
	}
	if len(p.Subject) == 0 {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_EmailSubjectIsNotSet, rme.Msg_EmailSubjectIsNotSet, nil)
	}
	if len(p.Message) == 0 {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_EmailMessageIsNotSet, rme.Msg_EmailMessageIsNotSet, nil)
	}

	c.GetMailerGuard().Lock()
	defer c.GetMailerGuard().Unlock()

	err := c.GetMailer().SendMail([]string{p.Recipient}, p.Subject, p.Message)
	if err != nil {
		c.logError(err)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_MailerError, fmt.Sprintf(rme.MsgF_MailerError, err.Error()), err)
	}

	result = &rm.SendEmailMessageResult{}

	return result, nil
}
