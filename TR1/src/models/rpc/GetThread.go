package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type GetThreadParams struct {
	CommonParams
	Thread *cm.Thread `json:"thread"`
}

type GetThreadResult struct {
	CommonResult
	Thread *cm.Thread `json:"thread"`
}
