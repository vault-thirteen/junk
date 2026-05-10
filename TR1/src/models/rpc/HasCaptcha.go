package rm

type HasCaptchaParams struct {
	TaskId string `json:"taskId"`
}

type HasCaptchaResult struct {
	CommonResult
	TaskId  string `json:"taskId"`
	IsFound bool   `json:"isFound"`
}
