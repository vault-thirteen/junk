package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type AddForumParams struct {
	CommonParams
	Forum *cm.Forum `json:"forum"`
}

type AddForumResult struct {
	CommonResult
	Forum *cm.Forum `json:"forum"`
}
