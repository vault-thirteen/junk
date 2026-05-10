package c

import (
	"bytes"
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) ConfirmLogOut(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ConfirmLogOutParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ConfirmLogOutResult
	r, re = c.confirmLogOut(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) confirmLogOut(p *rm.ConfirmLogOutParams) (result *rm.ConfirmLogOutResult, re *jrm1.RpcError) {
	var userWithSession *cm.User

	// Access check.
	{
		userWithSession, re = c.mustBeAnAuthToken(p.Auth)
		if re != nil {
			return nil, re
		}
	}

	// Check parameters.
	{
		if len(p.RequestId) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_RequestIdIsNotSet, rme.Msg_RequestIdIsNotSet, nil)
		}
		if !p.AreYouSure {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_NotSure, rme.Msg_NotSure, nil)
		}
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())
	var lor = &cm.LogOutRequest{
		RequestId: p.RequestId,
	}
	err = dbC.FindLogOutRequest(lor)
	if err != nil {
		return nil, c.databaseError(err)
	}

	// Check for fraud.
	{
		if !bytes.Equal(lor.UserIPAB, p.Auth.UserIPAB) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
		if lor.UserId != userWithSession.Id {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
	}

	// Confirm user logging out.
	{
		re = c.logOutUserBySelf(userWithSession.Id, userWithSession.Session.Id)
		if re != nil {
			return nil, re
		}

		err = dbC.DeleteLogOutRequest(lor)
		if err != nil {
			return nil, c.databaseError(err)
		}

		// Result.
		result = &rm.ConfirmLogOutResult{
			Success: rm.Success{OK: true},
		}

		return result, nil
	}
}
