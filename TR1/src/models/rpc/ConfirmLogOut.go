package rm

type ConfirmLogOutParams struct {
	CommonParams

	// Fields provided by user.
	RequestId  string `json:"requestId"`
	AreYouSure bool   `json:"areYouSure"`
}

type ConfirmLogOutResult struct {
	CommonResult
	Success
}
