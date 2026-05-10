package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
)

func (c *Controller) StartLogOut(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.StartLogOutParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.StartLogOutResult
	r, re = c.startLogOut(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) startLogOut(p *rm.StartLogOutParams) (result *rm.StartLogOutResult, re *jrm1.RpcError) {
	var userWithSession *cm.User

	// Access check.
	{
		userWithSession, re = c.mustBeAnAuthToken(p.Auth)
		if re != nil {
			return nil, re
		}
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())

	// Start logging out.
	{
		var requestId *string
		requestId, re = c.createRequestId()
		if re != nil {
			return nil, re
		}

		var lor = cm.LogOutRequest{
			UserId:    userWithSession.Id,
			RequestId: *requestId,
			UserIPAB:  p.Auth.UserIPAB,
		}
		err = dbC.CreateLogOutRequest(lor)
		if err != nil {
			return nil, c.databaseError(err)
		}

		result = &rm.StartLogOutResult{
			RequestId: *requestId,
		}

		return result, nil
	}
}
