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

func (c *Controller) ChangeMessageThread(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ChangeMessageThreadParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ChangeMessageThreadResult
	r, re = c.changeMessageThread(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) changeMessageThread(p *rm.ChangeMessageThreadParams) (result *rm.ChangeMessageThreadResult, re *jrm1.RpcError) {
	// Get caller user.
	var user *cm.User
	user, re = c.getSelfRoles(rm.GetSelfRolesParams{CommonParams: p.CommonParams})
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_AuthError, fmt.Sprintf(rme.MsgF_AuthError, re.AsError().Error()), nil)
	}

	// Access check.
	{
		if !user.Roles.IsAdministrator {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Permission, rme.Msg_Permission, nil)
		}
	}

	// Check parameters.
	{
		if p.Message == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_MessageIsNotSet, rme.Msg_MessageIsNotSet, nil)
		}
		if p.Message.Id == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
		if p.Message.ThreadId == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	err := dbC.ChangeMessageThread(user, p.Message)
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

	result = &rm.ChangeMessageThreadResult{
		Success: rm.Success{OK: true},
	}
	return result, nil
}
