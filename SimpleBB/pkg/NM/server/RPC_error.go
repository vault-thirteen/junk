package server

import "net/http"

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_UserIdIsNotSet            = 1
	RpcErrorCode_TextIsNotSet              = 2
	RpcErrorCode_NotificationIdIsNotSet    = 3
	RpcErrorCode_NotificationIsNotFound    = 4
	RpcErrorCode_NotificationIsAlreadyRead = 5
	RpcErrorCode_TestError                 = 6
	RpcErrorCode_PageIsNotSet              = 7
	RpcErrorCode_SystemEvent               = 8
	RpcErrorCode_ResourceIdIsNotSet        = 9
	RpcErrorCode_ResourceIsNotFound        = 10
	RpcErrorCode_ResourceIsNotValid        = 11
	RpcErrorCode_FormatStringType          = 12
	RpcErrorCode_ThreadIdIsNotSet          = 13
	RpcErrorCode_CreatorIsNotSet           = 14
)

// Messages.
const (
	RpcErrorMsg_UserIdIsNotSet            = "user ID is not set"
	RpcErrorMsg_TextIsNotSet              = "text is not set"
	RpcErrorMsg_NotificationIdIsNotSet    = "notification ID is not set"
	RpcErrorMsg_NotificationIsNotFound    = "notification is not found"
	RpcErrorMsg_NotificationIsAlreadyRead = "notification is already read"
	RpcErrorMsgF_TestError                = "test error: %s"
	RpcErrorMsg_PageIsNotSet              = "page is not set"
	RpcErrorMsg_SystemEvent               = "system event error"
	RpcErrorMsg_ResourceIdIsNotSet        = "resource ID is not set"
	RpcErrorMsg_ResourceIsNotFound        = "resource is not found"
	RpcErrorMsg_ResourceIsNotValid        = "resource is not valid"
	RpcErrorMsg_FormatStringType          = "format string type error"
	RpcErrorMsg_ThreadIdIsNotSet          = "thread ID is not set"
	RpcErrorMsg_CreatorIsNotSet           = "creator is not set"
)

// Unique HTTP status codes used in the map:
// - 400 (Bad request);
// - 404 (Not found);
// - 409 (Conflict);
// - 500 (Internal server error).
func GetMapOfHttpStatusCodesByRpcErrorCodes() map[int]int {
	return map[int]int{
		RpcErrorCode_UserIdIsNotSet:            http.StatusBadRequest,
		RpcErrorCode_TextIsNotSet:              http.StatusBadRequest,
		RpcErrorCode_NotificationIdIsNotSet:    http.StatusBadRequest,
		RpcErrorCode_NotificationIsNotFound:    http.StatusNotFound,
		RpcErrorCode_NotificationIsAlreadyRead: http.StatusConflict,
		RpcErrorCode_TestError:                 http.StatusInternalServerError,
		RpcErrorCode_PageIsNotSet:              http.StatusBadRequest,
		RpcErrorCode_SystemEvent:               http.StatusBadRequest,
		RpcErrorCode_ResourceIdIsNotSet:        http.StatusBadRequest,
		RpcErrorCode_ResourceIsNotFound:        http.StatusNotFound,
		RpcErrorCode_ResourceIsNotValid:        http.StatusBadRequest,
		RpcErrorCode_FormatStringType:          http.StatusBadRequest,
		RpcErrorCode_ThreadIdIsNotSet:          http.StatusBadRequest,
		RpcErrorCode_CreatorIsNotSet:           http.StatusBadRequest,
	}
}
