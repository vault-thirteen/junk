package c

import (
	"encoding/json"
	"fmt"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/RingCaptcha/models"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) CheckCaptcha(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.CheckCaptchaParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.CheckCaptchaResult
	r, re = c.checkCaptcha(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) checkCaptcha(p *rm.CheckCaptchaParams) (result *rm.CheckCaptchaResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.TaskId) == 0 {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaTaskIdIsNotSet, rme.Msg_CaptchaTaskIdIsNotSet, nil)
	}
	if p.Value == 0 {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsNotSet, rme.Msg_CaptchaAnswerIsNotSet, nil)
	}

	resp, err := c.far.cs.CheckCaptcha(&m.CheckCaptchaRequest{TaskId: p.TaskId, Value: p.Value})
	if err != nil {
		c.logError(err)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaCheckError, fmt.Sprintf(rme.Msg_CaptchaCheckError, err.Error()), nil)
	}

	result = &rm.CheckCaptchaResult{
		TaskId:    p.TaskId,
		IsSuccess: resp.IsSuccess,
	}

	return result, nil
}
