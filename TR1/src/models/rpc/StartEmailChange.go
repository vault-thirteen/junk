package rm

type StartEmailChangeParams struct {
	CommonParams

	// Fields provided by user.
	NewEmail string `json:"newEmail"`
}

type StartEmailChangeResult struct {
	CommonResult
	RequestId string `json:"requestId"`
	CaptchaId string `json:"captchaId"`
	AuthData  []byte `json:"authData"`
}
