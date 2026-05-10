package rm

type StartLogOutParams struct {
	CommonParams
}

type StartLogOutResult struct {
	CommonResult
	RequestId string `json:"requestId"`
}
