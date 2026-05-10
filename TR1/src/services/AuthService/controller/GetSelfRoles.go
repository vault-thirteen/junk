package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) GetSelfRoles(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.GetSelfRolesParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.GetSelfRolesResult
	r, re = c.getSelfRoles(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) getSelfRoles(p *rm.GetSelfRolesParams) (result *rm.GetSelfRolesResult, re *jrm1.RpcError) {
	var userWithSession *cm.User

	// Access check.
	{
		userWithSession, re = c.mustBeAnAuthToken(p.Auth)
		if re != nil {
			return nil, re
		}

		if !userWithSession.Roles.CanLogIn {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Permission, rme.Msg_Permission, nil)
		}
	}

	dbC := dbc.NewDbController(c.GetDb())

	user := &cm.User{Id: userWithSession.Id}
	err := dbC.GetUserRoles(user)
	if err != nil {
		return nil, c.databaseError(err)
	}

	c.attachUserSpecialRoles(user)

	result = &rm.GetSelfRolesResult{User: user}
	return result, nil
}
