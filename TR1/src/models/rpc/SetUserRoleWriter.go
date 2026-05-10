package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type SetUserRoleWriterParams struct {
	CommonParams
	User          *cm.User `json:"user"`
	IsRoleEnabled bool     `json:"isRoleEnabled"`
}

type SetUserRoleWriterResult struct {
	CommonResult
	Success
}
