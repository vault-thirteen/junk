package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type ApproveRegistrationRequestRFAParams struct {
	CommonParams
	User *cm.User `json:"user"`
}

type ApproveRegistrationRequestRFAResult struct {
	CommonResult
	Success
}
