package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type UnbanUserParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type UnbanUserResult struct {
	CommonResult
	Success
}
