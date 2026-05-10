package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type DeleteThreadParams struct {
	CommonParams
	Thread *cm.Thread `json:"thread"`
}

type DeleteThreadResult struct {
	CommonResult
	Success
}
