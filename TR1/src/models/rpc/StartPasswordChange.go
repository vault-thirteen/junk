package rm

type StartPasswordChangeParams struct {
	CommonParams

	// Fields provided by user.
	NewPassword  string `json:"newPassword"`
	NewPassword2 string `json:"newPassword2"`
}

type StartPasswordChangeResult struct {
	CommonResult
	RequestId string `json:"requestId"`
	CaptchaId string `json:"captchaId"`
	AuthData  []byte `json:"authData"`
}
