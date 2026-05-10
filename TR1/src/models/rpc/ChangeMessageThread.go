package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type ChangeMessageThreadParams struct {
	CommonParams
	Message *cm.Message `json:"message"`
}

type ChangeMessageThreadResult struct {
	CommonResult
	Success
}
