package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type MoveForumDownParams struct {
	CommonParams
	Forum *cm.Forum `json:"forum"`
}

type MoveForumDownResult struct {
	CommonResult
	Success
}
