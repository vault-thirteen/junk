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

func (c *Controller) AddThread(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.AddThreadParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.AddThreadResult
	r, re = c.addThread(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) addThread(p *rm.AddThreadParams) (result *rm.AddThreadResult, re *jrm1.RpcError) {
	// Get caller user.
	var user *cm.User
	user, re = c.getSelfRoles(rm.GetSelfRolesParams{CommonParams: p.CommonParams})
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_AuthError, fmt.Sprintf(rme.MsgF_AuthError, re.AsError().Error()), nil)
	}

	// Access check.
	{
		if !user.Roles.CanCreateThread {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Permission, rme.Msg_Permission, nil)
		}
	}

	// Check parameters.
	{
		if p.Forum == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_ForumIsNotSet, rme.Msg_ForumIsNotSet, nil)
		}
		if p.Forum.Id == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
		if p.Thread == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_ThreadIsNotSet, rme.Msg_ThreadIsNotSet, nil)
		}
		if len(p.Thread.Name) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_NameIsNotSet, rme.Msg_NameIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	thread, err := dbC.AddThread(p.Forum, p.Thread)
	if err != nil {
		return nil, c.databaseError(err)
	}

	result = &rm.AddThreadResult{
		Thread: thread,
	}
	return result, nil
}
