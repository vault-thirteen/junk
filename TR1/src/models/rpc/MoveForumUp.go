package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type MoveForumUpParams struct {
	CommonParams
	Forum *cm.Forum `json:"forum"`
}

type MoveForumUpResult struct {
	CommonResult
	Success
}
