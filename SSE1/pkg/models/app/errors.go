package app

// Errors Messages.
const (
	ErrAuthenticationFailure      = "Authentication Failure"
	ErrNullPointer                = "Null Pointer"
	ErrTokenDataMismatch          = "Token Data Mismatch"
	ErrTypeCastFailure            = "Type Cast Failure"
	ErrUnsupportedLoggerType      = "Unsupported Logger Type"
	ErrUnsupportedStorageType     = "Unsupported Storage Type"
	ErrUserCanNotBeDisabled       = "User can not be disabled"
	ErrCanNotLogOutOtherUsers     = "Can not log out other Users"
	ErrRegisteredUserDoesNotExist = "Registered User does not exist"
)

// Errors Formats.
const (
	ErrfBadRequestError               = "Bad Request: %v"
	ErrfCriticalError                 = "Critical Error: %v"
	ErrfForbiddenError                = "Forbidden: %v"
	ErrfHttpResponseError             = "HTTP Response Error: %v"
	ErrfUserAuthenticationNameIsTaken = "User Authentication Name is not available: '%v'."
)

// Sender Names.
const (
	SenderDisableUser                       = "DisableUser"
	SenderGetBrowserUserAgentId             = "GetBrowserUserAgentId"
	SenderGetHttpRequestBody                = "GetHttpRequestBody"
	SenderGetUserIdByAuthenticationName     = "GetUserIdByAuthenticationName"
	SenderHttpAuthentication                = "HTTP Authentication"
	SenderHttpHandlerApiUserDisable         = "httpHandlerApiUserDisable"
	SenderHttpHandlerApiUserLogOut          = "HttpHandlerApiUserLogOut"
	SenderHttpProtocolCheck                 = "HTTP Protocol Check"
	SenderIsUserAuthenticationNameFree      = "IsUserAuthenticationNameFree"
	SenderListRegisteredUsersPublicNames    = "ListRegisteredUsersPublicNames"
	SenderLogUserIn                         = "LogUserIn"
	SenderLogUserOut                        = "LogUserOut"
	SenderNewUserDisablingRequest           = "NewUserDisablingRequest"
	SenderNewUserLogInRequest               = "NewUserLogInRequest"
	SenderNewUserLogOutRequest              = "NewUserLogOutRequest"
	SenderNewUserRegistrationRequest        = "NewUserRegistrationRequest"
	SenderRegisteredUserIdExists            = "RegisteredUserIdExists"
	SenderRegisterUser                      = "RegisterUser"
	SenderRespondWithJsonObject             = "RespondWithJsonObject"
	SenderUnpackAuthDataFromContext         = "unpackAuthDataFromContext"
	SenderUpdateActiveSessionLastAccessTime = "UpdateActiveSessionLastAccessTime"
	SenderValidateMachineBrowserUserAgent   = "ValidateMachineBrowserUserAgent"
	SenderValidateMachineHost               = "ValidateMachineHost"
)
