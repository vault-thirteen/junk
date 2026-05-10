package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) ListUsers(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ListUsersParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ListUsersResult
	r, re = c.listUsers(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) listUsers(p *rm.ListUsersParams) (result *rm.ListUsersResult, re *jrm1.RpcError) {
	var userWithSession *cm.User

	// Access check.
	{
		userWithSession, re = c.mustBeAnAuthToken(p.Auth)
		if re != nil {
			return nil, re
		}

		if !userWithSession.Roles.IsAdministrator {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Permission, rme.Msg_Permission, nil)
		}
	}

	// Check parameters.
	{
		if p.Page == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PageIsNotSet, rme.Msg_PageIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbControllerWithPageSize(c.GetDb(), c.far.pageSize)

	users, totalCount, err := dbC.ListUsers(p.Page)
	if err != nil {
		return nil, c.databaseError(err)
	}

	c.attachUsersSpecialRoles(users)

	result = &rm.ListUsersResult{
		ItemsPaginated: rm.NewItemsPaginated[cm.User](p.Page, c.far.pageSize, users, totalCount),
	}
	return result, nil
}
