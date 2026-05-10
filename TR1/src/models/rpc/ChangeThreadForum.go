package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type ChangeThreadForumParams struct {
	CommonParams
	Thread *cm.Thread `json:"thread"`
}

type ChangeThreadForumResult struct {
	CommonResult
	Success
}
