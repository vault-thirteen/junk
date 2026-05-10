package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type ChangeThreadNameParams struct {
	CommonParams
	Thread *cm.Thread `json:"thread"`
}

type ChangeThreadNameResult struct {
	CommonResult
	Success
}
