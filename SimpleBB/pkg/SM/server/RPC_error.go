package server

import "net/http"

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_UserIdIsNotSet         = 1
	RpcErrorCode_ThreadIdIsNotSet       = 2
	RpcErrorCode_SubscriptionIsNotFound = 3
	RpcErrorCode_ThreadExists           = 4
	RpcErrorCode_ThreadDoesNotExist     = 5
	RpcErrorCode_TestError              = 6
	RpcErrorCode_PageIsNotSet           = 7
)

// Messages.
const (
	RpcErrorMsg_UserIdIsNotSet         = "user ID is not set"
	RpcErrorMsg_ThreadIdIsNotSet       = "thread ID is not set"
	RpcErrorMsg_SubscriptionIsNotFound = "subscription is not found"
	RpcErrorMsg_ThreadExists           = "thread exists"
	RpcErrorMsg_ThreadDoesNotExist     = "thread does not exist"
	RpcErrorMsgF_TestError             = "test error: %s"
	RpcErrorMsg_PageIsNotSet           = "page is not set"
)

// Unique HTTP status codes used in the map:
// - 400 (Bad request);
// - 404 (Not found);
// - 409 (Conflict);
// - 500 (Internal server error).
func GetMapOfHttpStatusCodesByRpcErrorCodes() map[int]int {
	return map[int]int{
		RpcErrorCode_UserIdIsNotSet:         http.StatusBadRequest,
		RpcErrorCode_ThreadIdIsNotSet:       http.StatusBadRequest,
		RpcErrorCode_SubscriptionIsNotFound: http.StatusNotFound,
		RpcErrorCode_ThreadExists:           http.StatusConflict,
		RpcErrorCode_ThreadDoesNotExist:     http.StatusConflict,
		RpcErrorCode_TestError:              http.StatusInternalServerError,
		RpcErrorCode_PageIsNotSet:           http.StatusBadRequest,
	}
}
