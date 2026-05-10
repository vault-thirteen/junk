package c

import (
	"encoding/json"
	"fmt"

	"github.com/vault-thirteen/JSON-RPC-M1"
	cm "github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) ChangeMessageText(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ChangeMessageTextParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ChangeMessageTextResult
	r, re = c.changeMessageText(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) changeMessageText(p *rm.ChangeMessageTextParams) (result *rm.ChangeMessageTextResult, re *jrm1.RpcError) {
	// Get caller user.
	var user *cm.User
	user, re = c.getSelfRoles(rm.GetSelfRolesParams{CommonParams: p.CommonParams})
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_AuthError, fmt.Sprintf(rme.MsgF_AuthError, re.AsError().Error()), nil)
	}

	// Check parameters.
	{
		if p.Message == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_MessageIsNotSet, rme.Msg_MessageIsNotSet, nil)
		}
		if p.Message.Id == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
		if len(p.Message.Text) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_TextIsNotSet, rme.Msg_TextIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	// Can this user change the message's text ?
	var canUserChangeMessageText bool
	canUserChangeMessageText, re = c.canUserChangeMessageText(user, p.Message)
	if re != nil {
		return nil, re
	}
	if !canUserChangeMessageText {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_UserCanNotChangeMessageText, rme.Msg_UserCanNotChangeMessageText, nil)
	}

	err := dbC.ChangeMessageText(user, p.Message)
	if err != nil {
		return nil, c.databaseError(err)
	}

	// Update the thread.
	{
		var message = &cm.Message{Id: p.Message.Id}
		message, err = dbC.GetMessage(message)
		if err != nil {
			return nil, c.databaseError(err)
		}

		thread := &cm.Thread{Id: message.ThreadId}
		err = dbC.TouchThread(thread)
		if err != nil {
			return nil, c.databaseError(err)
		}
	}

	result = &rm.ChangeMessageTextResult{
		Success: rm.Success{OK: true},
	}
	return result, nil
}
