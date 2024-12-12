package server

import "net/http"

// Common error codes start with 1024.

// Codes.
const (
	RpcErrorCode_FunctionIsNotImplemented = 1024 * 1
	RpcErrorCode_Authorisation            = 1024 * 2
	RpcErrorCode_Permission               = 1024 * 4
	RpcErrorCode_Database                 = 1024 * 8
	RpcErrorCode_RPCCall                  = 1024 * 16
	RpcErrorCode_UidList                  = 1024 * 32
	RpcErrorCode_Captcha                  = 1024 * 64
	RpcErrorCode_Password                 = 1024 * 128
	RpcErrorCode_ModuleSynchronisation    = 1024 * 256
	RpcErrorCode_SystemEvent              = 1024 * 512
)

// Messages.
const (
	RpcErrorMsg_FunctionIsNotImplemented = "function is not implemented"
	RpcErrorMsg_Authorisation            = "authorisation error"
	RpcErrorMsg_Permission               = "permission error"
	RpcErrorMsg_Database                 = "database error"
	RpcErrorMsg_RPCCall                  = "RPC call error"
	RpcErrorMsgF_UidList                 = "UidList error: %s" // Template.
	RpcErrorMsg_Captcha                  = "captcha error"
	RpcErrorMsg_Password                 = "password error"
	RpcErrorMsg_ModuleSynchronisation    = "module synchronisation error"
	RpcErrorMsg_SystemEvent              = "system event error"
)

// Unique HTTP status codes used in the map:
// - 403 (Forbidden);
// - 404 (Not found);
// - 500 (Internal server error).
func GetMapOfHttpStatusCodesByRpcErrorCodes() map[int]int {
	return map[int]int{
		RpcErrorCode_FunctionIsNotImplemented: http.StatusNotFound,
		RpcErrorCode_Authorisation:            http.StatusForbidden,
		RpcErrorCode_Permission:               http.StatusForbidden,
		RpcErrorCode_Database:                 http.StatusInternalServerError,
		RpcErrorCode_RPCCall:                  http.StatusInternalServerError,
		RpcErrorCode_UidList:                  http.StatusInternalServerError,
		RpcErrorCode_Captcha:                  http.StatusForbidden,
		RpcErrorCode_Password:                 http.StatusForbidden,
		RpcErrorCode_ModuleSynchronisation:    http.StatusInternalServerError,
		RpcErrorCode_SystemEvent:              http.StatusInternalServerError,
	}
}
