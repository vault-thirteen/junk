package c

import (
	"encoding/json"
	"fmt"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) ListMessages(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ListMessagesParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ListMessagesResult
	r, re = c.listMessages(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) listMessages(p *rm.ListMessagesParams) (result *rm.ListMessagesResult, re *jrm1.RpcError) {
	// Get caller user.
	var user *cm.User
	user, re = c.getSelfRoles(rm.GetSelfRolesParams{CommonParams: p.CommonParams})
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_AuthError, fmt.Sprintf(rme.MsgF_AuthError, re.AsError().Error()), nil)
	}

	// Access check.
	{
		if !user.Roles.CanRead {
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
		if p.Page == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PageIsNotSet, rme.Msg_PageIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	messages, totalCount, err := dbC.ListMessages(p.Thread, p.Page)
	if err != nil {
		return nil, c.databaseError(err)
	}

	result = &rm.ListMessagesResult{
		ItemsPaginated: rm.NewItemsPaginated[cm.Message](p.Page, c.far.pageSize, messages, totalCount),
	}
	return result, nil
}
