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

func (c *Controller) GetForum(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.GetForumParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.GetForumResult
	r, re = c.getForum(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) getForum(p *rm.GetForumParams) (result *rm.GetForumResult, re *jrm1.RpcError) {
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
		if p.Forum == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_ForumIsNotSet, rme.Msg_ForumIsNotSet, nil)
		}
		if p.Forum.Id == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	forum, err := dbC.GetForum(p.Forum)
	if err != nil {
		return nil, c.databaseError(err)
	}

	result = &rm.GetForumResult{
		Forum: forum,
	}
	return result, nil
}
