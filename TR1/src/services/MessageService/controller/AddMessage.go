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

func (c *Controller) AddMessage(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.AddMessageParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.AddMessageResult
	r, re = c.addMessage(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) addMessage(p *rm.AddMessageParams) (result *rm.AddMessageResult, re *jrm1.RpcError) {
	// Get caller user.
	var user *cm.User
	user, re = c.getSelfRoles(rm.GetSelfRolesParams{CommonParams: p.CommonParams})
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_AuthError, fmt.Sprintf(rme.MsgF_AuthError, re.AsError().Error()), nil)
	}

	// Access check.
	{
		if !user.Roles.CanWriteMessage {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Permission, rme.Msg_Permission, nil)
		}
	}

	// Check parameters.
	{
		if p.Thread == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_ThreadIsNotSet, rme.Msg_ThreadIsNotSet, nil)
		}
		if p.Thread.Id == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
		if p.Message == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_MessageIsNotSet, rme.Msg_MessageIsNotSet, nil)
		}
		if len(p.Message.Text) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_TextIsNotSet, rme.Msg_TextIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	// Can this user add a message ?
	var canUserAddMessage bool
	canUserAddMessage, re = c.canUserAddMessage(user, p.Thread)
	if re != nil {
		return nil, re
	}
	if !canUserAddMessage {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_UserCanNotAddMessage, rme.Msg_UserCanNotAddMessage, nil)
	}

	message, err := dbC.AddMessage(user, p.Thread, p.Message)
	if err != nil {
		return nil, c.databaseError(err)
	}

	// Update the thread.
	{
		err = dbC.TouchThread(p.Thread)
		if err != nil {
			return nil, c.databaseError(err)
		}
	}

	result = &rm.AddMessageResult{
		Message: message,
	}
	return result, nil
}
