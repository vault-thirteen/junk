package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) RejectRegistrationRequestRFA(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.RejectRegistrationRequestRFAParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.RejectRegistrationRequestRFAResult
	r, re = c.rejectRegistrationRequestRFA(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) rejectRegistrationRequestRFA(p *rm.RejectRegistrationRequestRFAParams) (result *rm.RejectRegistrationRequestRFAResult, re *jrm1.RpcError) {
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
		if len(p.User.Email) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_EmailIsNotSet, rme.Msg_EmailIsNotSet, nil)
		}
	}

	dbC := dbc.NewDbController(c.GetDb())

	rr := new(cm.RegistrationRequest)
	err := dbC.GetRegistrationRequestRFA(p.User.Email, rr)
	if err != nil {
		return nil, c.databaseError(err)
	}

	// Reject the request.
	err = dbC.DeleteRegistrationRequestRFA(rr)
	if err != nil {
		return nil, c.databaseError(err)
	}

	result = &rm.RejectRegistrationRequestRFAResult{
		Success: rm.Success{OK: true},
	}
	return result, nil
}
