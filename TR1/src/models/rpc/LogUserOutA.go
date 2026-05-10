package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type LogUserOutAParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type LogUserOutAResult struct {
	CommonResult
	Success
}
