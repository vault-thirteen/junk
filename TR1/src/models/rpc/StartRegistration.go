package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type StartRegistrationParams struct {
	CommonParams

	// Fields provided by user.
	User *cm.User `json:"user"`
}

type StartRegistrationResult struct {
	CommonResult
	RequestId string `json:"requestId"`
	CaptchaId string `json:"captchaId"`
}
