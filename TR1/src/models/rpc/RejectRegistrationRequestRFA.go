package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type RejectRegistrationRequestRFAParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type RejectRegistrationRequestRFAResult struct {
	CommonResult
	Success
}
