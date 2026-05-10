package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type BanUserParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type BanUserResult struct {
	CommonResult
	Success
}
