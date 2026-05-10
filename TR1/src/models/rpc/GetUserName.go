package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type GetUserNameParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type GetUserNameResult struct {
	CommonResult
	User *cm.User `json:"user"`
}
