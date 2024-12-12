package server

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_MailerError       = 1
	RpcErrorCode_RecipientIsNotSet = 2
	RpcErrorCode_SubjectIsNotSet   = 3
	RpcErrorCode_MessageIsNotSet   = 4
)

// Messages.
const (
	RpcErrorMsgF_MailerError      = "mailer error: %s"
	RpcErrorMsg_RecipientIsNotSet = "recipient is not set"
	RpcErrorMsg_SubjectIsNotSet   = "subject is not set"
	RpcErrorMsg_MessageIsNotSet   = "message is not set"
)
