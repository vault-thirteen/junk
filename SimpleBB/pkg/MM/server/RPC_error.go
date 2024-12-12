package server

import "net/http"

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_SectionNameIsNotSet      = 1
	RpcErrorCode_RootSectionAlreadyExists = 2
	RpcErrorCode_SectionIsNotFound        = 3
	RpcErrorCode_SectionIdIsNotSet        = 4
	RpcErrorCode_SectionHasChildren       = 5
	RpcErrorCode_RootSectionCanNotBeMoved = 6
	RpcErrorCode_ForumNameIsNotSet        = 7
	RpcErrorCode_ForumIsNotFound          = 8
	RpcErrorCode_ForumIdIsNotSet          = 9
	RpcErrorCode_ForumHasThreads          = 10
	RpcErrorCode_ThreadNameIsNotSet       = 11
	RpcErrorCode_ThreadIdIsNotSet         = 12
	RpcErrorCode_ThreadIsNotFound         = 13
	RpcErrorCode_ThreadIsNotEmpty         = 14
	RpcErrorCode_MessageTextIsNotSet      = 15
	RpcErrorCode_MessageIdIsNotSet        = 16
	RpcErrorCode_IncompatibleChildType    = 17
	RpcErrorCode_MessageIsNotFound        = 18
	RpcErrorCode_PageIsNotSet             = 19
	RpcErrorCode_TestError                = 20
)

// Messages.
const (
	RpcErrorMsg_SectionNameIsNotSet      = "section name is not set"
	RpcErrorMsg_RootSectionAlreadyExists = "root section already exists"
	RpcErrorMsg_SectionIsNotFound        = "section is not found"
	RpcErrorMsg_SectionIdIsNotSet        = "section ID is not set"
	RpcErrorMsg_SectionHasChildren       = "section has children"
	RpcErrorMsg_RootSectionCanNotBeMoved = "root section can not be moved"
	RpcErrorMsg_ForumNameIsNotSet        = "forum name is not set"
	RpcErrorMsg_ForumIsNotFound          = "forum is not found"
	RpcErrorMsg_ForumIdIsNotSet          = "forum ID is not set"
	RpcErrorMsg_ForumHasThreads          = "forum has threads"
	RpcErrorMsg_ThreadNameIsNotSet       = "thread name is not set"
	RpcErrorMsg_ThreadIdIsNotSet         = "thread ID is not set"
	RpcErrorMsg_ThreadIsNotFound         = "thread is not found"
	RpcErrorMsg_ThreadIsNotEmpty         = "thread is not empty"
	RpcErrorMsg_MessageTextIsNotSet      = "message text is not set"
	RpcErrorMsg_MessageIdIsNotSet        = "message ID is not set"
	RpcErrorMsg_IncompatibleChildType    = "incompatible child type"
	RpcErrorMsg_MessageIsNotFound        = "message is not found"
	RpcErrorMsg_PageIsNotSet             = "page is not set"
	RpcErrorMsgF_TestError               = "test error: %s"
)

// Unique HTTP status codes used in the map:
// - 400 (Bad request);
// - 404 (Not found);
// - 409 (Conflict);
// - 500 (Internal server error).
func GetMapOfHttpStatusCodesByRpcErrorCodes() map[int]int {
	return map[int]int{
		RpcErrorCode_SectionNameIsNotSet:      http.StatusBadRequest,
		RpcErrorCode_RootSectionAlreadyExists: http.StatusConflict,
		RpcErrorCode_SectionIsNotFound:        http.StatusNotFound,
		RpcErrorCode_SectionIdIsNotSet:        http.StatusBadRequest,
		RpcErrorCode_SectionHasChildren:       http.StatusConflict,
		RpcErrorCode_RootSectionCanNotBeMoved: http.StatusConflict,
		RpcErrorCode_ForumNameIsNotSet:        http.StatusBadRequest,
		RpcErrorCode_ForumIsNotFound:          http.StatusNotFound,
		RpcErrorCode_ForumIdIsNotSet:          http.StatusBadRequest,
		RpcErrorCode_ForumHasThreads:          http.StatusConflict,
		RpcErrorCode_ThreadNameIsNotSet:       http.StatusBadRequest,
		RpcErrorCode_ThreadIdIsNotSet:         http.StatusBadRequest,
		RpcErrorCode_ThreadIsNotFound:         http.StatusNotFound,
		RpcErrorCode_ThreadIsNotEmpty:         http.StatusConflict,
		RpcErrorCode_MessageTextIsNotSet:      http.StatusBadRequest,
		RpcErrorCode_MessageIdIsNotSet:        http.StatusBadRequest,
		RpcErrorCode_IncompatibleChildType:    http.StatusConflict,
		RpcErrorCode_MessageIsNotFound:        http.StatusNotFound,
		RpcErrorCode_PageIsNotSet:             http.StatusBadRequest,
		RpcErrorCode_TestError:                http.StatusInternalServerError,
	}
}
