package rm

type ConfirmPasswordChangeParams struct {
	CommonParams

	// Fields provided by user.
	RequestId        string `json:"requestId"`
	CaptchaAnswer    string `json:"captchaAnswer"`
	VerificationCode string `json:"verificationCode"`
	AuthData         []byte `json:"authData"`
}

type ConfirmPasswordChangeResult struct {
	CommonResult
	Success
	IsLoggedOut bool `json:"isLoggedOut"`
}
