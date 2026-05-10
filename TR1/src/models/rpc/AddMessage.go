package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type AddMessageParams struct {
	CommonParams
	Thread  *cm.Thread  `json:"thread"`
	Message *cm.Message `json:"message"`
}

type AddMessageResult struct {
	CommonResult
	Message *cm.Message `json:"message"`
}
