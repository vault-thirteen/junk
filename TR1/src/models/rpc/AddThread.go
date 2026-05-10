package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type AddThreadParams struct {
	CommonParams
	Forum  *cm.Forum  `json:"forum"`
	Thread *cm.Thread `json:"thread"`
}

type AddThreadResult struct {
	CommonResult
	Thread *cm.Thread `json:"thread"`
}
