package client

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// Message.
	FuncSendMessage = "SendMessage"

	// Other.
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
)
