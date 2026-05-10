package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/rpc"
)

func (c *Controller) Ping(_ *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	result = rm.PingResult{
		Success: rm.Success{OK: true},
	}
	return result, nil
}
