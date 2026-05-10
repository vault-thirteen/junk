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

func (c *Controller) AddForum(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.AddForumParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.AddForumResult
	r, re = c.addForum(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) addForum(p *rm.AddForumParams) (result *rm.AddForumResult, re *jrm1.RpcError) {
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
		if p.Forum == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_ForumIsNotSet, rme.Msg_ForumIsNotSet, nil)
		}
		if len(p.Forum.Name) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_NameIsNotSet, rme.Msg_NameIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	forum, err := dbC.AddForum(p.Forum)
	if err != nil {
		return nil, c.databaseError(err)
	}

	result = &rm.AddForumResult{
		Forum: forum,
	}
	return result, nil
}
