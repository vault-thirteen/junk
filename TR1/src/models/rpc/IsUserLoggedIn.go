package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type IsUserLoggedInParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type IsUserLoggedInResult struct {
	CommonResult
	IsUserLoggedIn bool `json:"isUserLoggedIn"`
}
