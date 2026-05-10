package rm

import cm "github.com/vault-thirteen/TR1/src/models/common"

type StartLogInParams struct {
	CommonParams

	// Fields provided by user.
	User *cm.User `json:"user"`
}

type StartLogInResult struct {
	CommonResult
	RequestId string `json:"requestId"`
	CaptchaId string `json:"captchaId"`
	AuthData  []byte `json:"authData"`
}
