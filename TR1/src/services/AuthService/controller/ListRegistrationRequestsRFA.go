package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) ListRegistrationRequestsRFA(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ListRegistrationRequestsRFAParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ListRegistrationRequestsRFAResult
	r, re = c.listRegistrationRequestsRFA(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) listRegistrationRequestsRFA(p *rm.ListRegistrationRequestsRFAParams) (result *rm.ListRegistrationRequestsRFAResult, re *jrm1.RpcError) {
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

	rrs, totalCount, err := dbC.ListRegistrationRequestsRFA(p.Page)
	if err != nil {
		return nil, c.databaseError(err)
	}

	result = &rm.ListRegistrationRequestsRFAResult{
		ItemsPaginated: rm.NewItemsPaginated[cm.RegistrationRequest](p.Page, c.far.pageSize, rrs, totalCount),
	}
	return result, nil
}
