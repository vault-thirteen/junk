package rm

type CreateCaptchaParams struct{}

type CreateCaptchaResult struct {
	CommonResult
	TaskId              string `json:"taskId"`
	ImageFormat         string `json:"imageFormat"`
	IsImageDataReturned bool   `json:"isImageDataReturned"`
	ImageDataB64        string `json:"imageDataB64,omitempty"`
}
