package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type ChangeMessageTextParams struct {
	CommonParams
	Message *cm.Message `json:"message"`
}

type ChangeMessageTextResult struct {
	CommonResult
	Success
}
