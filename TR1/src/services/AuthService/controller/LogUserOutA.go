package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) LogUserOutA(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.LogUserOutAParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.LogUserOutAResult
	r, re = c.logUserOutA(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) logUserOutA(p *rm.LogUserOutAParams) (result *rm.LogUserOutAResult, re *jrm1.RpcError) {
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
		if p.User == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserIsNotSet, rme.Msg_UserIsNotSet, nil)
		}
		if p.User.Id == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_IdIsNotSet, rme.Msg_IdIsNotSet, nil)
		}
		if p.User.Id == userWithSession.Id {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Permission, rme.Msg_Permission, nil)
		}
	}

	dbC := dbc.NewDbController(c.GetDb())

	userToLogOut := &cm.User{Id: p.User.Id}
	err := dbC.GetUserWithSessionByIdAbleToLogIn(userToLogOut)
	if err != nil {
		return nil, c.databaseError(err)
	}

	re = c.logOutUserByAdministrator(userToLogOut.Id, userToLogOut.Session.Id, &userWithSession.Id)
	if re != nil {
		return nil, re
	}

	result = &rm.LogUserOutAResult{Success: rm.Success{OK: true}}
	return result, nil
}
