package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type GetSelfRolesParams struct {
	CommonParams
}

type GetSelfRolesResult struct {
	CommonResult
	User *cm.User `json:"user"`
}
