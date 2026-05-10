package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type ListThreadsParams struct {
	CommonParams
	Forum *cm.Forum `json:"forum"`
	PageRequested
}

type ListThreadsResult struct {
	CommonResult
	ItemsPaginated
}
