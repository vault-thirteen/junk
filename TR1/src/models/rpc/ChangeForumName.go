package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type ChangeForumNameParams struct {
	CommonParams
	Forum *cm.Forum `json:"forum"`
}

type ChangeForumNameResult struct {
	CommonResult
	Success
}
