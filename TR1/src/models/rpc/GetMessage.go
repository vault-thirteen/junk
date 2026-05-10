package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type GetMessageParams struct {
	CommonParams
	Message *cm.Message `json:"message"`
}

type GetMessageResult struct {
	CommonResult
	Message *cm.Message `json:"message"`
}
