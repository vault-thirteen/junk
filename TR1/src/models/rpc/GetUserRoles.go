package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type GetUserRolesParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type GetUserRolesResult struct {
	CommonResult
	User *cm.User `json:"user"`
}
