package rm

const (
	RpcDurationFieldName  = "dur"
	RpcRequestIdFieldName = "rid"
)

const (
	ServiceNextPingAttemptDelaySec     = 5
	ServicePingAttemptsDurationMinutes = 15
)

const (
	UrlSchemeHttp  = "http"
	UrlSchemeHttps = "https"
)

const (
	Func_Ping = "Ping"
)

const (
	ServiceName_AuthService    = "AuthService"
	ServiceName_CaptchaService = "CaptchaService"
	ServiceName_GatewayService = "GatewayService"
	ServiceName_MailerService  = "MailerService"
	ServiceName_MessageService = "MessageService"
)

const (
	ServiceShortName_Auth    = "auth"
	ServiceShortName_Captcha = "captcha" // Captcha images (Proxy).
	ServiceShortName_Gateway = "gateway"
	ServiceShortName_Mailer  = "mailer"
	ServiceShortName_Message = "message"
	ServiceShortName_RCS     = "rcs" // Captcha questions (RPC).
)

const (
	ServiceId_Auth    = 1
	ServiceId_Captcha = 2
	ServiceId_Gateway = 3
	ServiceId_Mailer  = 4
	ServiceId_Message = 5
	ServiceId_RCS     = 6
)

// Names of internal RPC functions.

// Auth service.
const (
	Func_ApproveRegistrationRequestRFA = "ApproveRegistrationRequestRFA"
	Func_BanUser                       = "BanUser"
	Func_ConfirmEmailChange            = "ConfirmEmailChange"
	Func_ConfirmLogIn                  = "ConfirmLogIn"
	Func_ConfirmLogOut                 = "ConfirmLogOut"
	Func_ConfirmPasswordChange         = "ConfirmPasswordChange"
	Func_ConfirmRegistration           = "ConfirmRegistration"
	Func_GetSelfRoles                  = "GetSelfRoles"
	Func_GetUserName                   = "GetUserName"
	Func_GetUserParameters             = "GetUserParameters"
	Func_GetUserRoles                  = "GetUserRoles"
	Func_GetUserSession                = "GetUserSession"
	Func_IsUserLoggedIn                = "IsUserLoggedIn"
	Func_ListRegistrationRequestsRFA   = "ListRegistrationRequestsRFA"
	Func_ListUsers                     = "ListUsers"
	Func_ListUserSessions              = "ListUserSessions"
	Func_LogUserOutA                   = "LogUserOutA"
	Func_RejectRegistrationRequestRFA  = "RejectRegistrationRequestRFA"
	Func_SetUserRoleAuthor             = "SetUserRoleAuthor"
	Func_SetUserRoleReader             = "SetUserRoleReader"
	Func_SetUserRoleWriter             = "SetUserRoleWriter"
	Func_StartEmailChange              = "StartEmailChange"
	Func_StartLogIn                    = "StartLogIn"
	Func_StartLogOut                   = "StartLogOut"
	Func_StartPasswordChange           = "StartPasswordChange"
	Func_StartRegistration             = "StartRegistration"
	Func_UnbanUser                     = "UnbanUser"
)

// Captcha service.
const (
	Func_CreateCaptcha = "CreateCaptcha"
	Func_CheckCaptcha  = "CheckCaptcha"
	Func_HasCaptcha    = "HasCaptcha"
)

// Mailer service.
const (
	Func_SendEmailMessage = "SendEmailMessage"
)

// Message service.
const (
	Func_AddForum            = "AddForum"
	Func_AddMessage          = "AddMessage"
	Func_AddThread           = "AddThread"
	Func_ChangeForumName     = "ChangeForumName"
	Func_ChangeMessageText   = "ChangeMessageText"
	Func_ChangeMessageThread = "ChangeMessageThread"
	Func_ChangeThreadForum   = "ChangeThreadForum"
	Func_ChangeThreadName    = "ChangeThreadName"
	Func_DeleteForum         = "DeleteForum"
	Func_DeleteMessage       = "DeleteMessage"
	Func_DeleteThread        = "DeleteThread"
	Func_GetForum            = "GetForum"
	Func_GetMessage          = "GetMessage"
	Func_GetThread           = "GetThread"
	Func_ListForums          = "ListForums"
	Func_ListMessages        = "ListMessages"
	Func_ListThreads         = "ListThreads"
	Func_MoveForumDown       = "MoveForumDown"
	Func_MoveForumUp         = "MoveForumUp"
)

// Names of external (gateway) RPC functions.
const (
	// AuthService.
	ApiFunctionName_ApproveRegistrationRequestRFA = "approveRegistrationRequestRFA"
	ApiFunctionName_BanUser                       = "banUser"
	ApiFunctionName_ConfirmEmailChange            = "confirmEmailChange"
	ApiFunctionName_ConfirmLogIn                  = "confirmLogIn"
	ApiFunctionName_ConfirmLogOut                 = "confirmLogOut"
	ApiFunctionName_ConfirmPasswordChange         = "confirmPasswordChange"
	ApiFunctionName_ConfirmRegistration           = "confirmRegistration"
	ApiFunctionName_GetSelfRoles                  = "getSelfRoles"
	ApiFunctionName_GetUserName                   = "getUserName"
	ApiFunctionName_GetUserParameters             = "getUserParameters"
	ApiFunctionName_GetUserRoles                  = "getUserRoles"
	ApiFunctionName_GetUserSession                = "getUserSession"
	ApiFunctionName_IsUserLoggedIn                = "isUserLoggedIn"
	ApiFunctionName_ListRegistrationRequestsRFA   = "listRegistrationRequestsRFA"
	ApiFunctionName_ListUsers                     = "listUsers"
	ApiFunctionName_ListUserSessions              = "listUserSessions"
	ApiFunctionName_LogUserOutA                   = "logUserOutA"
	ApiFunctionName_RejectRegistrationRequestRFA  = "rejectRegistrationRequestRFA"
	ApiFunctionName_SetUserRoleAuthor             = "setUserRoleAuthor"
	ApiFunctionName_SetUserRoleReader             = "setUserRoleReader"
	ApiFunctionName_SetUserRoleWriter             = "setUserRoleWriter"
	ApiFunctionName_StartEmailChange              = "startEmailChange"
	ApiFunctionName_StartLogIn                    = "startLogIn"
	ApiFunctionName_StartLogOut                   = "startLogOut"
	ApiFunctionName_StartPasswordChange           = "startPasswordChange"
	ApiFunctionName_StartRegistration             = "startRegistration"
	ApiFunctionName_UnbanUser                     = "unbanUser"

	// MessageService.
	ApiFunctionName_AddForum            = "addForum"
	ApiFunctionName_AddMessage          = "addMessage"
	ApiFunctionName_AddThread           = "addThread"
	ApiFunctionName_ChangeForumName     = "changeForumName"
	ApiFunctionName_ChangeMessageText   = "changeMessageText"
	ApiFunctionName_ChangeMessageThread = "changeMessageThread"
	ApiFunctionName_ChangeThreadForum   = "changeThreadForum"
	ApiFunctionName_ChangeThreadName    = "changeThreadName"
	ApiFunctionName_DeleteForum         = "deleteForum"
	ApiFunctionName_DeleteMessage       = "deleteMessage"
	ApiFunctionName_DeleteThread        = "deleteThread"
	ApiFunctionName_GetForum            = "getForum"
	ApiFunctionName_GetMessage          = "getMessage"
	ApiFunctionName_GetThread           = "getThread"
	ApiFunctionName_ListForums          = "listForums"
	ApiFunctionName_ListMessages        = "listMessages"
	ApiFunctionName_ListThreads         = "listThreads"
	ApiFunctionName_MoveForumDown       = "moveForumDown"
	ApiFunctionName_MoveForumUp         = "moveForumUp"
)
