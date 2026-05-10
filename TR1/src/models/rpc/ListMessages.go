package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type ListMessagesParams struct {
	CommonParams
	Thread *cm.Thread `json:"thread"`
	PageRequested
}

type ListMessagesResult struct {
	CommonResult
	ItemsPaginated
}
