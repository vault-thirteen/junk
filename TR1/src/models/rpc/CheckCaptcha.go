package rm

type CheckCaptchaParams struct {
	TaskId string `json:"taskId"`
	Value  uint   `json:"value"`
}

type CheckCaptchaResult struct {
	CommonResult
	TaskId    string `json:"taskId"`
	IsSuccess bool   `json:"isSuccess"`
}
