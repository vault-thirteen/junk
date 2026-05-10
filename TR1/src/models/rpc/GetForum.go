package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type GetForumParams struct {
	CommonParams
	Forum *cm.Forum `json:"forum"`
}

type GetForumResult struct {
	CommonResult
	Forum *cm.Forum `json:"forum"`
}
