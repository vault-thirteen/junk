package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type DeleteForumParams struct {
	CommonParams
	Forum *cm.Forum `json:"forum"`
}

type DeleteForumResult struct {
	CommonResult
	Success
}
