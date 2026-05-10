package rm

type ConfirmLogInParams struct {
	CommonParams

	// Fields provided by user.
	RequestId        string `json:"requestId"`
	CaptchaAnswer    string `json:"captchaAnswer"`
	VerificationCode string `json:"verificationCode"`
	AuthData         []byte `json:"authData"`
}

type ConfirmLogInResult struct {
	CommonResult
	Success
	IsTokenSet bool   `json:"isTokenSet"`
	Token      string `json:"token"`
}
