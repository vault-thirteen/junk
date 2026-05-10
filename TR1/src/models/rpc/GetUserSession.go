package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type GetUserSessionParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type GetUserSessionResult struct {
	CommonResult
	Session *cm.Session `json:"session"`
}
