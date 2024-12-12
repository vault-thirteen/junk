package client

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// Captcha.
	FuncCreateCaptcha = "CreateCaptcha"
	FuncCheckCaptcha  = "CheckCaptcha"

	// Other.
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
)
