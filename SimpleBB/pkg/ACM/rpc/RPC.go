package rpc

import (
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

// Ping.

type PingParams = rpc2.PingParams
type PingResult = rpc2.PingResult

// User registration.

type RegisterUserParams struct {
	rpc2.CommonParams

	// Step number.
	StepN simple.StepNumber `json:"stepN"`

	// E-mail address. Is used on steps 1, 2 and 3.
	Email simple.Email `json:"email"`

	// Verification code. Is used on steps 2 and 3.
	VerificationCode simple.VerificationCode `json:"verificationCode"`

	// Name. Is used on step 3.
	Name simple.Name `json:"name"`

	// Password. Is used on step 3.
	Password simple.Password `json:"password"`
}
type RegisterUserResult struct {
	rpc2.CommonResult

	// Next required step. If set to zero, no further step is required.
	NextStep simple.StepNumber `json:"nextStep"`
}

type GetListOfRegistrationsReadyForApprovalParams struct {
	rpc2.CommonParams
	Page base2.Count `json:"page"`
}
type GetListOfRegistrationsReadyForApprovalResult struct {
	rpc2.CommonResult
	RRFA     []models.RegistrationReadyForApproval `json:"rrfa"`
	PageData *rpc2.PageData                        `json:"pageData,omitempty"`
}

type RejectRegistrationRequestParams struct {
	rpc2.CommonParams
	RegistrationRequestId base2.Id `json:"registrationRequestId"`
}
type RejectRegistrationRequestResult = rpc2.CommonResultWithSuccess

type ApproveAndRegisterUserParams struct {
	rpc2.CommonParams

	// E-mail address of an approved user.
	Email simple.Email `json:"email"`
}
type ApproveAndRegisterUserResult = rpc2.CommonResultWithSuccess

// Logging in and out.

type LogUserInParams struct {
	rpc2.CommonParams

	// Step number.
	StepN simple.StepNumber `json:"stepN"`

	// E-mail address.
	// This is the main identifier of a user.
	// It is used on all steps.
	Email simple.Email `json:"email"`

	// Request ID.
	// It protects preliminary sessions from being hi-jacked.
	// Is used on steps 2 and 3.
	RequestId simple.RequestId `json:"requestId"`

	// Captcha answer.
	// This field is optional and may be used on step 2.
	CaptchaAnswer simple.CaptchaAnswer `json:"captchaAnswer"`

	// Authentication data provided for the challenge.
	// Is used on step 2.
	AuthChallengeResponse rpc2.AuthChallengeData `json:"authChallengeResponse"`

	// Verification Code.
	// Is used on step 3.
	VerificationCode simple.VerificationCode `json:"verificationCode"`
}
type LogUserInResult struct {
	rpc2.CommonResult

	// Next required step. If set to zero, no further step is required.
	NextStep simple.StepNumber `json:"nextStep"`

	RequestId     simple.RequestId       `json:"requestId"`
	AuthDataBytes rpc2.AuthChallengeData `json:"authDataBytes"`

	// Captcha parameters.
	IsCaptchaNeeded base2.Flag        `json:"isCaptchaNeeded"`
	CaptchaId       *simple.CaptchaId `json:"captchaId"`

	// JWT key maker.
	IsWebTokenSet  base2.Flag            `json:"isWebTokenSet"`
	WebTokenString simple.WebTokenString `json:"wts,omitempty"`
}

type LogUserOutParams struct {
	rpc2.CommonParams
}
type LogUserOutResult = rpc2.CommonResultWithSuccess

type LogUserOutAParams struct {
	rpc2.CommonParams
	UserId base2.Id `json:"userId"`
}
type LogUserOutAResult = rpc2.CommonResultWithSuccess

type GetListOfLoggedUsersParams struct {
	rpc2.CommonParams
}
type GetListOfLoggedUsersResult struct {
	rpc2.CommonResult
	LoggedUserIds []base2.Id `json:"loggedUserIds"`
}

type GetListOfLoggedUsersOnPageParams struct {
	rpc2.CommonParams
	Page base2.Count `json:"page"`
}
type GetListOfLoggedUsersOnPageResult struct {
	rpc2.CommonResult
	LoggedUserIds []base2.Id     `json:"loggedUserIds"`
	PageData      *rpc2.PageData `json:"pageData,omitempty"`
}

type GetListOfAllUsersParams struct {
	rpc2.CommonParams
}
type GetListOfAllUsersResult struct {
	rpc2.CommonResult
	UserIds []base2.Id `json:"userIds"`
}

type GetListOfAllUsersOnPageParams struct {
	rpc2.CommonParams
	Page base2.Count `json:"page"`
}
type GetListOfAllUsersOnPageResult struct {
	rpc2.CommonResult
	UserIds  []base2.Id     `json:"userIds"`
	PageData *rpc2.PageData `json:"pageData,omitempty"`
}

type IsUserLoggedInParams struct {
	rpc2.CommonParams
	UserId base2.Id `json:"userId"`
}
type IsUserLoggedInResult struct {
	rpc2.CommonResult
	UserId         base2.Id   `json:"userId"`
	IsUserLoggedIn base2.Flag `json:"isUserLoggedIn"`
}

// Various actions.

type ChangePasswordParams struct {
	rpc2.CommonParams

	// Step number.
	StepN simple.StepNumber `json:"stepN"`

	// New password.
	// Is used on step 1.
	NewPassword simple.Password `json:"newPassword"`

	// Request ID.
	// It protects password changes from being hi-jacked.
	// Is used on step 2.
	RequestId simple.RequestId `json:"requestId"`

	// Authentication data provided for the challenge.
	// Is used on step 2.
	AuthChallengeResponse rpc2.AuthChallengeData `json:"authChallengeResponse"`

	// Verification Code.
	// Is used on step 2.
	VerificationCode simple.VerificationCode `json:"verificationCode"`

	// Captcha answer.
	// This field is optional and may be used on step 2.
	CaptchaAnswer simple.CaptchaAnswer `json:"captchaAnswer"`
}
type ChangePasswordResult struct {
	rpc2.CommonResult
	rpc2.Success

	// Next required step. If set to zero, no further step is required.
	NextStep simple.StepNumber `json:"nextStep"`

	RequestId     simple.RequestId       `json:"requestId"`
	AuthDataBytes rpc2.AuthChallengeData `json:"authDataBytes"`

	// Captcha parameters.
	IsCaptchaNeeded base2.Flag       `json:"isCaptchaNeeded"`
	CaptchaId       simple.CaptchaId `json:"captchaId"`
}

type ChangeEmailParams struct {
	rpc2.CommonParams

	// Step number.
	StepN simple.StepNumber `json:"stepN"`

	// New e-mail address.
	// Is used on step 1.
	NewEmail simple.Email `json:"newEmail"`

	// Request ID.
	// It protects e-mail changes from being hi-jacked.
	// Is used on step 2.
	RequestId simple.RequestId `json:"requestId"`

	// Authentication data provided for the challenge.
	// Is used on step 2.
	AuthChallengeResponse rpc2.AuthChallengeData `json:"authChallengeResponse"`

	// Verification Code for the old e-mail.
	// Is used on step 2.
	VerificationCodeOld simple.VerificationCode `json:"verificationCodeOld"`

	// Verification Code for the new e-mail.
	// Is used on step 2.
	VerificationCodeNew simple.VerificationCode `json:"verificationCodeNew"`

	// Captcha answer.
	// This field is optional and may be used on step 2.
	CaptchaAnswer simple.CaptchaAnswer `json:"captchaAnswer"`
}
type ChangeEmailResult struct {
	rpc2.CommonResult
	rpc2.Success

	// Next required step. If set to zero, no further step is required.
	NextStep simple.StepNumber `json:"nextStep"`

	RequestId     simple.RequestId       `json:"requestId"`
	AuthDataBytes rpc2.AuthChallengeData `json:"authDataBytes"`

	// Captcha parameters.
	IsCaptchaNeeded base2.Flag       `json:"isCaptchaNeeded"`
	CaptchaId       simple.CaptchaId `json:"captchaId"`
}

type GetUserSessionParams struct {
	rpc2.CommonParams
	UserId base2.Id `json:"userId"`
}
type GetUserSessionResult struct {
	rpc2.CommonResult
	User    derived1.IUser  `json:"user"`
	Session *models.Session `json:"session"`
}

// User properties.

type GetUserNameParams struct {
	rpc2.CommonParams
	UserId base2.Id `json:"userId"`
}
type GetUserNameResult struct {
	rpc2.CommonResult
	User derived1.IUser `json:"user"`
}

type GetUserRolesParams struct {
	rpc2.CommonParams
	UserId base2.Id `json:"userId"`
}
type GetUserRolesResult struct {
	rpc2.CommonResult
	User derived1.IUser `json:"user"`
}

type ViewUserParametersParams struct {
	rpc2.CommonParams
	UserId base2.Id `json:"userId"`
}
type ViewUserParametersResult struct {
	rpc2.CommonResult
	User derived1.IUser `json:"user"`
}

type SetUserRoleAuthorParams = SetUserRoleCommonParams
type SetUserRoleAuthorResult = SetUserRoleCommonResult

type SetUserRoleWriterParams = SetUserRoleCommonParams
type SetUserRoleWriterResult = SetUserRoleCommonResult

type SetUserRoleReaderParams = SetUserRoleCommonParams
type SetUserRoleReaderResult = SetUserRoleCommonResult

type GetSelfRolesParams struct {
	rpc2.CommonParams
}
type GetSelfRolesResult struct {
	rpc2.CommonResult
	User derived1.IUser `json:"user"`
}

// User banning.

type BanUserParams struct {
	rpc2.CommonParams
	UserId base2.Id `json:"userId"`
}
type BanUserResult = rpc2.CommonResultWithSuccess

type UnbanUserParams struct {
	rpc2.CommonParams
	UserId base2.Id `json:"userId"`
}
type UnbanUserResult = rpc2.CommonResultWithSuccess

// Other.

type ShowDiagnosticDataParams struct{}
type ShowDiagnosticDataResult struct {
	rpc2.CommonResult
	rpc2.RequestsCount
}

type TestParams struct{}
type TestResult struct {
	rpc2.CommonResult
}

// Common models.

type SetUserRoleCommonParams struct {
	rpc2.CommonParams
	UserId        base2.Id   `json:"userId"`
	IsRoleEnabled base2.Flag `json:"isRoleEnabled"`
}
type SetUserRoleCommonResult = rpc2.CommonResultWithSuccess
