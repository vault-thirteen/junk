package c

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) CreateCaptcha(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.CreateCaptchaParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.CreateCaptchaResult
	r, re = c.createCaptcha(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) createCaptcha(p *rm.CreateCaptchaParams) (result *rm.CreateCaptchaResult, re *jrm1.RpcError) {
	ccResponse, err := c.far.cs.CreateCaptcha()
	if err != nil {
		c.logError(err)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaCreationError, fmt.Sprintf(rme.MsgF_CaptchaCreationError, err.Error()), nil)
	}

	result = &rm.CreateCaptchaResult{
		TaskId:              ccResponse.TaskId,
		ImageFormat:         ccResponse.ImageFormat,
		IsImageDataReturned: ccResponse.IsImageDataReturned,
	}

	if ccResponse.IsImageDataReturned {
		result.ImageDataB64 = base64.StdEncoding.EncodeToString(ccResponse.ImageData)
	}

	return result, nil
}
