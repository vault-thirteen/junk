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

func (c *Controller) ChangeThreadForum(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ChangeThreadForumParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ChangeThreadForumResult
	r, re = c.changeThreadForum(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) changeThreadForum(p *rm.ChangeThreadForumParams) (result *rm.ChangeThreadForumResult, re *jrm1.RpcError) {
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
		if p.Thread == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_ThreadIsNotSet, rme.Msg_ThreadIsNotSet, nil)
		}
		if p.Thread.Id == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
		if p.Thread.ForumId == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	err := dbC.ChangeThreadForum(p.Thread)
	if err != nil {
		return nil, c.databaseError(err)
	}

	result = &rm.ChangeThreadForumResult{
		Success: rm.Success{OK: true},
	}
	return result, nil
}
