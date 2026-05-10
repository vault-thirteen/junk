package rm

type ConfirmRegistrationParams struct {
	CommonParams

	// Fields provided by user.
	RequestId        string `json:"requestId"`
	CaptchaAnswer    string `json:"captchaAnswer"`
	VerificationCode string `json:"verificationCode"`
}

type ConfirmRegistrationResult struct {
	CommonResult
	Success
	IsApprovalRequired bool `json:"isApprovalRequired"`
}
