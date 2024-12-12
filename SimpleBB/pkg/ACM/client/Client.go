package c

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// User registration.
	FuncRegisterUser                           = "RegisterUser"
	FuncGetListOfRegistrationsReadyForApproval = "GetListOfRegistrationsReadyForApproval"
	FuncRejectRegistrationRequest              = "RejectRegistrationRequest"
	FuncApproveAndRegisterUser                 = "ApproveAndRegisterUser"

	// Logging in and out.
	FuncLogUserIn                  = "LogUserIn"
	FuncLogUserOut                 = "LogUserOut"
	FuncLogUserOutA                = "LogUserOutA"
	FuncGetListOfLoggedUsers       = "GetListOfLoggedUsers"
	FuncGetListOfLoggedUsersOnPage = "GetListOfLoggedUsersOnPage"
	FuncGetListOfAllUsers          = "GetListOfAllUsers"
	FuncGetListOfAllUsersOnPage    = "GetListOfAllUsersOnPage"
	FuncIsUserLoggedIn             = "IsUserLoggedIn"

	// Various actions.
	FuncChangePassword = "ChangePassword"
	FuncChangeEmail    = "ChangeEmail"
	FuncGetUserSession = "GetUserSession"

	// User properties.
	FuncGetUserName        = "GetUserName"
	FuncGetUserRoles       = "GetUserRoles"
	FuncViewUserParameters = "ViewUserParameters"
	FuncSetUserRoleAuthor  = "SetUserRoleAuthor"
	FuncSetUserRoleWriter  = "SetUserRoleWriter"
	FuncSetUserRoleReader  = "SetUserRoleReader"
	FuncGetSelfRoles       = "GetSelfRoles"

	// User banning.
	FuncBanUser   = "BanUser"
	FuncUnbanUser = "UnbanUser"

	// Other.
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
	FuncTest               = "Test"
)
