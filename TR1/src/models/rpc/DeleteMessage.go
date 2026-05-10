package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type DeleteMessageParams struct {
	CommonParams
	Message *cm.Message `json:"message"`
}

type DeleteMessageResult struct {
	CommonResult
	Success
}
