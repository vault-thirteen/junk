package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type GetUserParametersParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type GetUserParametersResult struct {
	CommonResult
	User *cm.User `json:"user"`
}
