package server

import (
	"net/http"
)

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_AuthChallengeResponseIsNotSet      = 1
	RpcErrorCode_BPP_GenerateRandomSalt             = 2
	RpcErrorCode_BPP_PackSymbols                    = 3
	RpcErrorCode_CaptchaAnswerIsNotSet              = 4
	RpcErrorCode_CaptchaAnswerIsWrong               = 5
	RpcErrorCode_EmailAddressIsNotValid             = 6
	RpcErrorCode_EmailAddressIsUsed                 = 7
	RpcErrorCode_EmailChangeIsNotFound              = 8
	RpcErrorCode_JWTCreation                        = 9
	RpcErrorCode_NameIsNotSet                       = 10
	RpcErrorCode_NameIsTooLong                      = 11
	RpcErrorCode_NameIsUsed                         = 12
	RpcErrorCode_NewPasswordIsNotSet                = 13
	RpcErrorCode_NewEmailIsNotSet                   = 14
	RpcErrorCode_PasswordChangeIsNotFound           = 15
	RpcErrorCode_PasswordIsNotValid                 = 16
	RpcErrorCode_PasswordIsNotSet                   = 17
	RpcErrorCode_PasswordIsTooLong                  = 18
	RpcErrorCode_PasswordIsWrong                    = 19
	RpcErrorCode_RequestIdGenerator                 = 20
	RpcErrorCode_RequestIdIsNotSet                  = 21
	RpcErrorCode_SmtpModule                         = 22
	RpcErrorCode_StepIsUnknown                      = 23
	RpcErrorCode_UserAlreadyStartedToChangePassword = 24
	RpcErrorCode_UserAlreadyStartedToChangeEmail    = 25
	RpcErrorCode_UserCanNotLogIn                    = 26
	RpcErrorCode_UserHasAlreadyStartedToLogIn       = 27
	RpcErrorCode_UserHasNotStartedToLogIn           = 28
	RpcErrorCode_UserIdIsNotSet                     = 29
	RpcErrorCode_UserIsAlreadyLoggedIn              = 30
	RpcErrorCode_UserIsNotFound                     = 31
	RpcErrorCode_UserPreSessionIsNotFound           = 32
	RpcErrorCode_VerificationCodeGenerator          = 33
	RpcErrorCode_VerificationCodeIsNotSet           = 34
	RpcErrorCode_VerificationCodeIsWrong            = 35
	RpcErrorCode_PageIsNotSet                       = 36
	RpcErrorCode_DatabaseInconsistency              = 37
	RpcErrorCode_SessionIsNotFound                  = 38
	RpcErrorCode_UserNameIsNotFound                 = 39
	RpcErrorCode_EmailAddressIsNotSet               = 40
	RpcErrorCode_CaptchaIdIsNotSet                  = 41
)

// Messages.
const (
	RpcErrorMsg_AuthChallengeResponseIsNotSet      = "authorisation challenge response is not set"
	RpcErrorMsg_BPP_GenerateRandomSalt             = "BPP: GenerateRandomSalt error"
	RpcErrorMsgF_BPP_PackSymbols                   = "BPP: PackSymbols error: %s" // Template.
	RpcErrorMsg_CaptchaAnswerIsNotSet              = "captcha answer is not set"
	RpcErrorMsg_CaptchaAnswerIsWrong               = "captcha answer is wrong"
	RpcErrorMsg_EmailAddressIsNotValid             = "e-mail address is not valid"
	RpcErrorMsg_EmailAddressIsUsed                 = "e-mail address is used"
	RpcErrorMsg_EmailChangeIsNotFound              = "request for e-mail address change is not found"
	RpcErrorMsgF_JWTCreation                       = "JWT creation error: %s" // Template.
	RpcErrorMsg_NameIsNotSet                       = "name is not set"
	RpcErrorMsg_NameIsTooLong                      = "name is too long"
	RpcErrorMsg_NameIsUsed                         = "name is already used"
	RpcErrorMsg_NewPasswordIsNotSet                = "new password is not set"
	RpcErrorMsg_NewEmailIsNotSet                   = "new e-mail address is not set"
	RpcErrorMsg_PasswordChangeIsNotFound           = "request for password change is not found"
	RpcErrorMsg_PasswordIsNotValid                 = "password is not valid"
	RpcErrorMsg_PasswordIsNotSet                   = "password is not set"
	RpcErrorMsg_PasswordIsTooLong                  = "password is too long"
	RpcErrorMsg_PasswordIsWrong                    = "password is wrong"
	RpcErrorMsg_RequestIdGenerator                 = "error generating request ID"
	RpcErrorMsg_RequestIdIsNotSet                  = "request ID is not set"
	RpcErrorMsg_SmtpModule                         = "SMTP module error"
	RpcErrorMsg_StepIsUnknown                      = "unknown step"
	RpcErrorMsg_UserAlreadyStartedToChangePassword = "user has already started to change password"
	RpcErrorMsg_UserAlreadyStartedToChangeEmail    = "user has already started to change e-mail address"
	RpcErrorMsg_UserCanNotLogIn                    = "user can not log in"
	RpcErrorMsg_UserHasAlreadyStartedToLogIn       = "user has already started to log in"
	RpcErrorMsg_UserHasNotStartedToLogIn           = "user has not started to log in"
	RpcErrorMsg_UserIdIsNotSet                     = "user ID is not set"
	RpcErrorMsg_UserIsAlreadyLoggedIn              = "user is already logged in"
	RpcErrorMsg_UserIsNotFound                     = "user is not found"
	RpcErrorMsg_UserPreSessionIsNotFound           = "user's preliminary session is not found"
	RpcErrorMsg_VerificationCodeGenerator          = "verification code generator error"
	RpcErrorMsg_VerificationCodeIsNotSet           = "verification code is not set"
	RpcErrorMsg_VerificationCodeIsWrong            = "verification code is wrong"
	RpcErrorMsg_PageIsNotSet                       = "page is not set"
	RpcErrorMsg_DatabaseInconsistency              = "database inconsistency"
	RpcErrorMsg_SessionIsNotFound                  = "session is not found"
	RpcErrorMsg_UserNameIsNotFound                 = "user name is not found"
	RpcErrorMsg_EmailAddressIsNotSet               = "email address is not set"
	RpcErrorMsg_CaptchaIdIsNotSet                  = "captcha ID is not set"
)

// Unique HTTP status codes used in the map:
// - 400 (Bad request);
// - 403 (Forbidden);
// - 404 (Not found);
// - 409 (Conflict);
// - 500 (Internal server error).
func GetMapOfHttpStatusCodesByRpcErrorCodes() map[int]int {
	return map[int]int{
		RpcErrorCode_AuthChallengeResponseIsNotSet:      http.StatusBadRequest,
		RpcErrorCode_BPP_GenerateRandomSalt:             http.StatusInternalServerError,
		RpcErrorCode_BPP_PackSymbols:                    http.StatusInternalServerError,
		RpcErrorCode_CaptchaAnswerIsNotSet:              http.StatusBadRequest,
		RpcErrorCode_CaptchaAnswerIsWrong:               http.StatusForbidden,
		RpcErrorCode_EmailAddressIsNotValid:             http.StatusBadRequest,
		RpcErrorCode_EmailAddressIsUsed:                 http.StatusConflict,
		RpcErrorCode_EmailChangeIsNotFound:              http.StatusNotFound,
		RpcErrorCode_JWTCreation:                        http.StatusInternalServerError,
		RpcErrorCode_NameIsNotSet:                       http.StatusBadRequest,
		RpcErrorCode_NameIsTooLong:                      http.StatusBadRequest,
		RpcErrorCode_NameIsUsed:                         http.StatusConflict,
		RpcErrorCode_NewPasswordIsNotSet:                http.StatusBadRequest,
		RpcErrorCode_NewEmailIsNotSet:                   http.StatusBadRequest,
		RpcErrorCode_PasswordChangeIsNotFound:           http.StatusNotFound,
		RpcErrorCode_PasswordIsNotValid:                 http.StatusBadRequest,
		RpcErrorCode_PasswordIsNotSet:                   http.StatusBadRequest,
		RpcErrorCode_PasswordIsTooLong:                  http.StatusBadRequest,
		RpcErrorCode_PasswordIsWrong:                    http.StatusForbidden,
		RpcErrorCode_RequestIdGenerator:                 http.StatusInternalServerError,
		RpcErrorCode_RequestIdIsNotSet:                  http.StatusBadRequest,
		RpcErrorCode_SmtpModule:                         http.StatusInternalServerError,
		RpcErrorCode_StepIsUnknown:                      http.StatusBadRequest,
		RpcErrorCode_UserAlreadyStartedToChangePassword: http.StatusForbidden,
		RpcErrorCode_UserAlreadyStartedToChangeEmail:    http.StatusForbidden,
		RpcErrorCode_UserCanNotLogIn:                    http.StatusForbidden,
		RpcErrorCode_UserHasAlreadyStartedToLogIn:       http.StatusForbidden,
		RpcErrorCode_UserHasNotStartedToLogIn:           http.StatusForbidden,
		RpcErrorCode_UserIdIsNotSet:                     http.StatusBadRequest,
		RpcErrorCode_UserIsAlreadyLoggedIn:              http.StatusForbidden,
		RpcErrorCode_UserIsNotFound:                     http.StatusNotFound,
		RpcErrorCode_UserPreSessionIsNotFound:           http.StatusNotFound,
		RpcErrorCode_VerificationCodeGenerator:          http.StatusInternalServerError,
		RpcErrorCode_VerificationCodeIsNotSet:           http.StatusBadRequest,
		RpcErrorCode_VerificationCodeIsWrong:            http.StatusForbidden,
		RpcErrorCode_PageIsNotSet:                       http.StatusBadRequest,
		RpcErrorCode_DatabaseInconsistency:              http.StatusInternalServerError,
		RpcErrorCode_SessionIsNotFound:                  http.StatusNotFound,
		RpcErrorCode_UserNameIsNotFound:                 http.StatusNotFound,
		RpcErrorCode_EmailAddressIsNotSet:               http.StatusBadRequest,
		RpcErrorCode_CaptchaIdIsNotSet:                  http.StatusInternalServerError,
	}
}
