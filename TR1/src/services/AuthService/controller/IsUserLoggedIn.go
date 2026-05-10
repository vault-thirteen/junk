package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) IsUserLoggedIn(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.IsUserLoggedInParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.IsUserLoggedInResult
	r, re = c.isUserLoggedIn(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) isUserLoggedIn(p *rm.IsUserLoggedInParams) (result *rm.IsUserLoggedInResult, re *jrm1.RpcError) {
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

	// Check parameters.
	{
		if p.User == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserIsNotSet, rme.Msg_UserIsNotSet, nil)
		}
		if p.User.Id == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbController(c.GetDb())

	user := &cm.User{Id: p.User.Id}
	result = new(rm.IsUserLoggedInResult)

	_, err := dbC.FindUserSession(user)
	if err != nil {
		result.IsUserLoggedIn = false
	} else {
		result.IsUserLoggedIn = true
	}

	return result, nil
}
