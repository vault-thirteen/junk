package rpc

import (
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// Ping.

type PingParams = rpc2.PingParams
type PingResult = rpc2.PingResult

// Captcha.

type CreateCaptchaParams struct{}
type CreateCaptchaResult struct {
	rpc2.CommonResult
	TaskId              string `json:"taskId"`
	ImageFormat         string `json:"imageFormat"`
	IsImageDataReturned bool   `json:"isImageDataReturned"`
	ImageDataB64        string `json:"imageDataB64,omitempty"`
}

type CheckCaptchaParams struct {
	TaskId string `json:"taskId"`
	Value  uint   `json:"value"`
}
type CheckCaptchaResult struct {
	rpc2.CommonResult
	TaskId    string `json:"taskId"`
	IsSuccess bool   `json:"isSuccess"`
}

// Other.

type ShowDiagnosticDataParams struct{}
type ShowDiagnosticDataResult struct {
	rpc2.CommonResult
	rpc2.RequestsCount
}
