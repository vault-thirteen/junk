package server

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_FirewallIsDisabled = 1
	RpcErrorCode_IPAddressIsNotSet  = 2
	RpcErrorCode_BlockTimeIsNotSet  = 3
)

// Messages.
const (
	RpcErrorMsg_FirewallIsDisabled = "Firewall is disabled"
	RpcErrorMsg_IPAddressIsNotSet  = "IP address is not set"
	RpcErrorMsg_BlockTimeIsNotSet  = "Block time is not set"
)
