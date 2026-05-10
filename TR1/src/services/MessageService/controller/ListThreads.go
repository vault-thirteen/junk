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

func (c *Controller) ListThreads(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ListThreadsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ListThreadsResult
	r, re = c.listThreads(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) listThreads(p *rm.ListThreadsParams) (result *rm.ListThreadsResult, re *jrm1.RpcError) {
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
		if p.Page == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PageIsNotSet, rme.Msg_PageIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	threads, totalCount, err := dbC.ListThreads(p.Forum, p.Page)
	if err != nil {
		return nil, c.databaseError(err)
	}

	result = &rm.ListThreadsResult{
		ItemsPaginated: rm.NewItemsPaginated[cm.Thread](p.Page, c.far.pageSize, threads, totalCount),
	}
	return result, nil
}
