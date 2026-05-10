package rm

type ConfirmEmailChangeParams struct {
	CommonParams

	// Fields provided by user.
	RequestId         string `json:"requestId"`
	CaptchaAnswer     string `json:"captchaAnswer"`
	VerificationCodeA string `json:"verificationCodeA"`
	VerificationCodeB string `json:"verificationCodeB"`
	AuthData          []byte `json:"authData"`
}

type ConfirmEmailChangeResult struct {
	CommonResult
	Success
	IsLoggedOut bool `json:"isLoggedOut"`
}
