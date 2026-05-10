package c

import (
	"encoding/json"
	"fmt"

	"github.com/vault-thirteen/JSON-RPC-M1"
	rcm "github.com/vault-thirteen/RingCaptcha/models"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) HasCaptcha(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.HasCaptchaParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.HasCaptchaResult
	r, re = c.hasCaptcha(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) hasCaptcha(p *rm.HasCaptchaParams) (result *rm.HasCaptchaResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.TaskId) == 0 {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaTaskIdIsNotSet, rme.Msg_CaptchaTaskIdIsNotSet, nil)
	}

	resp, err := c.far.cs.HasCaptcha(&rcm.HasCaptchaRequest{TaskId: p.TaskId})
	if err != nil {
		c.logError(err)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaError, fmt.Sprintf(rme.MsgF_CaptchaError, err.Error()), nil)
	}

	result = &rm.HasCaptchaResult{
		TaskId:  p.TaskId,
		IsFound: resp.IsFound,
	}

	return result, nil
}
