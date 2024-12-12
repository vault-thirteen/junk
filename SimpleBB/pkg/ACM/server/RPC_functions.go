package server

import (
	"fmt"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/ACM/rpc"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/rpc"
	base22 "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UserRoles"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/IncidentType"
	le "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/LogEvent"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/LogEventType"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/User"
	rpc3 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"time"

	bpp "github.com/vault-thirteen/BytePackedPassword"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
)

// RPC functions.

// User registration.

func (srv *Server) registerUser(p *rpc2.RegisterUserParams) (result *rpc2.RegisterUserResult, re *jrm1.RpcError) {
	re = srv.mustBeNoAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	switch p.StepN {
	case 1:
		return srv.registerUserStep1(p)
	case 2:
		return srv.registerUserStep2(p)
	case 3:
		return srv.registerUserStep3(p)
	default:
		// Step is not supported.
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_StepIsUnknown, RpcErrorMsg_StepIsUnknown, nil)
	}
}

func (srv *Server) registerUserStep1(p *rpc2.RegisterUserParams) (result *rpc2.RegisterUserResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Email) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsNotSet, RpcErrorMsg_EmailAddressIsNotSet, nil)
	}

	// Is e-mail address valid ?
	re = srv.isEmailOfUserValid(p.Email)
	if re != nil {
		return nil, re
	}

	// Check that ient's e-mail address is not used.
	usersCount, err := srv.dbo.CountUsersWithEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	isAddrFree := usersCount == 0
	if !isAddrFree {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsUsed, RpcErrorMsg_EmailAddressIsUsed, nil)
	}

	// Add the client to a list of pre-registered users.
	// Database should contain restrictions against the following duplicate
	// properties of a user and a pre-user: name, email.
	err = srv.dbo.InsertPreRegisteredUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Create a verification code.
	var verificationCode *simple.VerificationCode
	verificationCode, re = srv.createVerificationCode()
	if re != nil {
		return nil, re
	}

	// Attach the verification code to the user.
	err = srv.dbo.AttachVerificationCodeToPreRegUser(p.Email, *verificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Send an e-mail message with a verification code.
	re = srv.sendVerificationCodeForReg(p.Email, *verificationCode)
	if re != nil {
		return nil, re
	}

	// Confirm the email message send.
	err = srv.dbo.SetPreRegUserEmailSendStatus(true, p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &rpc2.RegisterUserResult{NextStep: 2}, nil
}

func (srv *Server) registerUserStep2(p *rpc2.RegisterUserParams) (result *rpc2.RegisterUserResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Email) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsNotSet, RpcErrorMsg_EmailAddressIsNotSet, nil)
	}

	// Is e-mail address valid ?
	re = srv.isEmailOfUserValid(p.Email)
	if re != nil {
		return nil, re
	}

	// Is verification code set ?
	if len(p.VerificationCode) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsNotSet, RpcErrorMsg_VerificationCodeIsNotSet, nil)
	}

	// Check the verification code.
	ok, err := srv.dbo.CheckVerificationCodeForPreReg(p.Email, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		// Verification code can not be guessed.
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_VerificationCodeMismatch), p.Email, p.Auth.UserIPAB)

		// Delete the pre-registered user.
		err = srv.dbo.DeletePreRegUserIfNotApprovedByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsWrong, RpcErrorMsg_VerificationCodeIsWrong, nil)
	}

	// Mark the client as verified by e-mail address.
	err = srv.dbo.ApproveUserByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &rpc2.RegisterUserResult{NextStep: 3}, nil
}

func (srv *Server) registerUserStep3(p *rpc2.RegisterUserParams) (result *rpc2.RegisterUserResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Email) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsNotSet, RpcErrorMsg_EmailAddressIsNotSet, nil)
	}

	// Is e-mail address valid ?
	re = srv.isEmailOfUserValid(p.Email)
	if re != nil {
		return nil, re
	}

	// Is verification code set ?
	if len(p.VerificationCode) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsNotSet, RpcErrorMsg_VerificationCodeIsNotSet, nil)
	}

	// Is name set ?
	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NameIsNotSet, RpcErrorMsg_NameIsNotSet, nil)
	}

	// Is password set ?
	if len(p.Password) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordIsNotSet, RpcErrorMsg_PasswordIsNotSet, nil)
	}

	// Check name.
	if len([]byte(p.Name)) > int(srv.settings.SystemSettings.UserNameMaxLenInBytes) {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NameIsTooLong, RpcErrorMsg_NameIsTooLong, nil)
	}

	usersCount, err := srv.dbo.CountUsersWithName(p.Name)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	isUserNameFree := usersCount == 0
	if !isUserNameFree {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NameIsUsed, RpcErrorMsg_NameIsUsed, nil)
	}

	// Check password.
	re = isPasswordAllowed(p.Password)
	if re != nil {
		return nil, re
	}

	var pwdBytes []byte
	pwdBytes, err = bpp.PackSymbols([]rune(p.Password))
	if err != nil {
		// This error is very unlikely to happen.
		// So if it occurs, then it is an anomaly.
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_BPP_PackSymbols, fmt.Sprintf(RpcErrorMsgF_BPP_PackSymbols, err.Error()), nil)
	}

	if len(pwdBytes) > int(srv.settings.SystemSettings.UserPasswordMaxLenInBytes) {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordIsTooLong, RpcErrorMsg_PasswordIsTooLong, nil)
	}

	// Check the verification code.
	var ok bool
	ok, err = srv.dbo.CheckVerificationCodeForPreReg(p.Email, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_VerificationCodeMismatch), p.Email, p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsWrong, RpcErrorMsg_VerificationCodeIsWrong, nil)
	}

	// Set user's name and password.
	// Database should contain restrictions against the following duplicate
	// properties of a user and a pre-user: name, email.
	err = srv.dbo.SetPreRegUserData(p.Email, p.VerificationCode, p.Name, pwdBytes)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if srv.settings.SystemSettings.IsAdminApprovalRequired {
		return &rpc2.RegisterUserResult{NextStep: 4}, nil
	}

	err = srv.dbo.ApprovePreRegUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.RegisterPreRegUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	re = srv.sendGreetingAfterReg(p.Email)
	if re != nil {
		return nil, re
	}

	return &rpc2.RegisterUserResult{NextStep: 0}, nil
}

func (srv *Server) getListOfRegistrationsReadyForApproval(p *rpc2.GetListOfRegistrationsReadyForApprovalParams) (result *rpc2.GetListOfRegistrationsReadyForApprovalResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	rrfaCount, err := srv.dbo.CountRegistrationsReadyForApproval()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var rrfas []am.RegistrationReadyForApproval
	rrfas, err = srv.dbo.GetListOfRegistrationsReadyForApproval(p.Page, srv.settings.SystemSettings.PageSize)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.GetListOfRegistrationsReadyForApprovalResult{
		PageData: &rpc3.PageData{
			PageNumber:  p.Page,
			TotalPages:  base2.CalculateTotalPages(rrfaCount, srv.settings.SystemSettings.PageSize),
			PageSize:    srv.settings.SystemSettings.PageSize,
			ItemsOnPage: base2.Count(len(rrfas)),
			TotalItems:  rrfaCount,
		},
		RRFA: rrfas,
	}

	return result, nil
}

func (srv *Server) rejectRegistrationRequest(p *rpc2.RejectRegistrationRequestParams) (result *rpc2.RejectRegistrationRequestResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.RegistrationRequestId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RequestIdIsNotSet, RpcErrorMsg_RequestIdIsNotSet, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	err := srv.dbo.RejectRegistrationRequest(p.RegistrationRequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.RejectRegistrationRequestResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

func (srv *Server) approveAndRegisterUser(p *rpc2.ApproveAndRegisterUserParams) (result *rpc2.ApproveAndRegisterUserResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Email) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsNotSet, RpcErrorMsg_EmailAddressIsNotSet, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), p.Email, p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Is e-mail address valid ?
	re = srv.isEmailOfUserValid(p.Email)
	if re != nil {
		return nil, re
	}

	err := srv.dbo.ApprovePreRegUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.RegisterPreRegUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	re = srv.sendGreetingAfterReg(p.Email)
	if re != nil {
		return nil, re
	}

	result = &rpc2.ApproveAndRegisterUserResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// Logging in and out.

func (srv *Server) logUserIn(p *rpc2.LogUserInParams) (result *rpc2.LogUserInResult, re *jrm1.RpcError) {
	re = srv.mustBeNoAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	switch p.StepN {
	case 1:
		return srv.logUserInStep1(p)
	case 2:
		return srv.logUserInStep2(p)
	case 3:
		return srv.logUserInStep3(p)
	default:
		// Step is not supported.
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_StepIsUnknown, RpcErrorMsg_StepIsUnknown, nil)
	}
}

func (srv *Server) logUserInStep1(p *rpc2.LogUserInParams) (result *rpc2.LogUserInResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Email) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsNotSet, RpcErrorMsg_EmailAddressIsNotSet, nil)
	}

	// Is e-mail address valid ?
	re = srv.isEmailOfUserValid(p.Email)
	if re != nil {
		return nil, re
	}

	usersCount, err := srv.dbo.CountUsersWithEmailAbleToLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	canUserLogIn := usersCount == 1
	if !canUserLogIn {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserCanNotLogIn, RpcErrorMsg_UserCanNotLogIn, nil)
	}

	err = srv.dbo.DeleteAbandonedPreSessions()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var sessionsCount base2.Count
	sessionsCount, err = srv.dbo.CountSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	areUserSessionsPresent := sessionsCount > 0
	if areUserSessionsPresent {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_DoubleLogInAttempt), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIsAlreadyLoggedIn, RpcErrorMsg_UserIsAlreadyLoggedIn, nil)
	}

	var preSessionsCount base2.Count
	preSessionsCount, err = srv.dbo.CountPreSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	areUserPreSessionsPresent := preSessionsCount > 0
	if areUserPreSessionsPresent {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PreSessionHacking), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserHasAlreadyStartedToLogIn, RpcErrorMsg_UserHasAlreadyStartedToLogIn, nil)
	}

	var userId base2.Id
	userId, err = srv.dbo.GetUserIdByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var requestId *simple.RequestId
	requestId, re = srv.createRequestIdForLogIn()
	if re != nil {
		return nil, re
	}

	var pwdSalt []byte
	pwdSalt, err = bpp.GenerateRandomSalt()
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_BPP_GenerateRandomSalt, RpcErrorMsg_BPP_GenerateRandomSalt, nil)
	}

	// Create a captcha if needed.
	var isCaptchaNeeded base2.Flag
	isCaptchaNeeded, err = srv.isCaptchaNeededForLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var captchaData *rm.CreateCaptchaResult
	var captchaId *simple.CaptchaId
	if isCaptchaNeeded {
		captchaData, re = srv.createCaptcha()
		if re != nil {
			return nil, re
		}

		captchaId = (*simple.CaptchaId)(&captchaData.TaskId)
	}

	err = srv.dbo.CreatePreSession(userId, *requestId, p.Auth.UserIPAB, pwdSalt, isCaptchaNeeded, captchaId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.LogUserInResult{NextStep: 2, RequestId: *requestId, AuthDataBytes: pwdSalt}

	if isCaptchaNeeded {
		result.CaptchaId = captchaId
	} else {
		result.IsCaptchaNeeded = false
	}

	return result, nil
}

func (srv *Server) logUserInStep2(p *rpc2.LogUserInParams) (result *rpc2.LogUserInResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Email) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsNotSet, RpcErrorMsg_EmailAddressIsNotSet, nil)
	}

	// Is e-mail address valid ?
	re = srv.isEmailOfUserValid(p.Email)
	if re != nil {
		return nil, re
	}

	if len(p.RequestId) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RequestIdIsNotSet, RpcErrorMsg_RequestIdIsNotSet, nil)
	}

	if len(p.AuthChallengeResponse) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_AuthChallengeResponseIsNotSet, RpcErrorMsg_AuthChallengeResponseIsNotSet, nil)
	}

	usersCount, err := srv.dbo.CountUsersWithEmailAbleToLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	canUserLogIn := usersCount == 1
	if !canUserLogIn {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserCanNotLogIn, RpcErrorMsg_UserCanNotLogIn, nil)
	}

	err = srv.dbo.DeleteAbandonedPreSessions()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var sessionsCount base2.Count
	sessionsCount, err = srv.dbo.CountSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	areUserSessionsPresent := sessionsCount > 0
	if areUserSessionsPresent {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_DoubleLogInAttempt), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIsAlreadyLoggedIn, RpcErrorMsg_UserIsAlreadyLoggedIn, nil)
	}

	var preSessionsCount base2.Count
	preSessionsCount, err = srv.dbo.CountPreSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if preSessionsCount != 1 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PreSessionHacking), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserHasNotStartedToLogIn, RpcErrorMsg_UserHasNotStartedToLogIn, nil)
	}

	var userId base2.Id
	userId, err = srv.dbo.GetUserIdByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Get the request to log in.
	var preSession *am.PreSession
	preSession, err = srv.dbo.GetPreSessionByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if preSession == nil {
		// Request Id code can not be guessed.
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PreSessionHacking), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserPreSessionIsNotFound, RpcErrorMsg_UserPreSessionIsNotFound, nil)
	}

	// Check the Request ID (indirectly).
	if preSession.UserId != userId {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PreSessionHacking), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserPreSessionIsNotFound, RpcErrorMsg_UserPreSessionIsNotFound, nil)
	}

	// Check the IP address.
	if !p.Auth.UserIPAB.Equal(preSession.UserIPAB) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PreSessionHacking), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserPreSessionIsNotFound, RpcErrorMsg_UserPreSessionIsNotFound, nil)
	}

	// Check the captcha answer if needed.
	var ccr *rm.CheckCaptchaResult
	if preSession.IsCaptchaRequired {
		if len(p.CaptchaAnswer) == 0 {
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CaptchaAnswerIsNotSet, RpcErrorMsg_CaptchaAnswerIsNotSet, nil)
		}
		if preSession.CaptchaId == nil {
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CaptchaIdIsNotSet, RpcErrorMsg_CaptchaIdIsNotSet, nil)
		}

		// Check captcha answer.
		ccr, re = srv.checkCaptcha(*preSession.CaptchaId, p.CaptchaAnswer)
		if re != nil {
			return nil, re
		}

		if !ccr.IsSuccess {
			srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_CaptchaAnswerMismatch), p.Email, p.Auth.UserIPAB)

			err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			// When captcha guess is wrong, we delete the preliminary session
			// to start the process from the first step.
			err = srv.dbo.DeletePreSessionByRequestId(preSession.RequestId)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CaptchaAnswerIsWrong, RpcErrorMsg_CaptchaAnswerIsWrong, nil)
		}
	}

	// Check the password.
	var ok bool
	ok, re = srv.checkPassword(userId, preSession.AuthDataBytes, p.AuthChallengeResponse)
	if re != nil {
		return nil, re
	}

	if !ok {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PasswordMismatch), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// When password is wrong, we delete the preliminary session to start
		// the process from the first step.
		err = srv.dbo.DeletePreSessionByRequestId(preSession.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordIsWrong, RpcErrorMsg_PasswordIsWrong, nil)
	}

	// Set flags.
	//
	// Captcha flag is set to true when either it was not required or it was
	// required and the answer was correct.
	err = srv.dbo.SetPreSessionCaptchaFlags(userId, preSession.RequestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.SetPreSessionPasswordFlag(userId, preSession.RequestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Create a new Request ID for the next step.
	var step3requestId *simple.RequestId
	step3requestId, re = srv.createRequestIdForLogIn()
	if re != nil {
		return nil, re
	}

	err = srv.dbo.UpdatePreSessionRequestId(userId, preSession.RequestId, *step3requestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Verification by E-mail.
	var verificationCode *simple.VerificationCode
	verificationCode, re = srv.createVerificationCode()
	if re != nil {
		return nil, re
	}

	err = srv.dbo.AttachVerificationCodeToPreSession(userId, *step3requestId, *verificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	re = srv.sendVerificationCodeForLogIn(p.Email, *verificationCode)
	if re != nil {
		return nil, re
	}

	err = srv.dbo.SetPreSessionEmailSendStatus(userId, *step3requestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.LogUserInResult{NextStep: 3, RequestId: *step3requestId}

	return result, nil
}

func (srv *Server) logUserInStep3(p *rpc2.LogUserInParams) (result *rpc2.LogUserInResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Email) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsNotSet, RpcErrorMsg_EmailAddressIsNotSet, nil)
	}

	// Is e-mail address valid ?
	re = srv.isEmailOfUserValid(p.Email)
	if re != nil {
		return nil, re
	}

	if len(p.RequestId) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RequestIdIsNotSet, RpcErrorMsg_RequestIdIsNotSet, nil)
	}

	if len(p.VerificationCode) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsNotSet, RpcErrorMsg_VerificationCodeIsNotSet, nil)
	}

	usersCount, err := srv.dbo.CountUsersWithEmailAbleToLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	canUserLogIn := usersCount == 1
	if !canUserLogIn {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserCanNotLogIn, RpcErrorMsg_UserCanNotLogIn, nil)
	}

	err = srv.dbo.DeleteAbandonedPreSessions()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var sessionsCount base2.Count
	sessionsCount, err = srv.dbo.CountSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	areUserSessionsPresent := sessionsCount > 0
	if areUserSessionsPresent {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_DoubleLogInAttempt), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIsAlreadyLoggedIn, RpcErrorMsg_UserIsAlreadyLoggedIn, nil)
	}

	var preSessionsCount base2.Count
	preSessionsCount, err = srv.dbo.CountPreSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if preSessionsCount != 1 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PreSessionHacking), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserHasNotStartedToLogIn, RpcErrorMsg_UserHasNotStartedToLogIn, nil)
	}

	var userId base2.Id
	userId, err = srv.dbo.GetUserIdByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Get the request to log in.
	var preSession *am.PreSession
	preSession, err = srv.dbo.GetPreSessionByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if preSession == nil {
		// Request Id code can not be guessed.
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PreSessionHacking), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserPreSessionIsNotFound, RpcErrorMsg_UserPreSessionIsNotFound, nil)
	}

	// Check the Request ID (indirectly).
	if preSession.UserId != userId {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PreSessionHacking), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserPreSessionIsNotFound, RpcErrorMsg_UserPreSessionIsNotFound, nil)
	}

	// Check the IP address.
	if !p.Auth.UserIPAB.Equal(preSession.UserIPAB) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PreSessionHacking), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserPreSessionIsNotFound, RpcErrorMsg_UserPreSessionIsNotFound, nil)
	}

	// Check the verification code.
	var ok bool
	ok, err = srv.dbo.CheckVerificationCodeForLogIn(p.RequestId, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		// Verification code can not be guessed.
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_VerificationCodeMismatch), p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// Delete the pre-session on error to avoid brute force checks.
		err = srv.dbo.DeletePreSessionByRequestId(preSession.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsWrong, RpcErrorMsg_VerificationCodeIsWrong, nil)
	}

	// Set verification flags.
	err = srv.dbo.SetPreSessionVerificationFlag(userId, preSession.RequestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Create a normal session and delete the preliminary one.
	var sessionId base2.Id
	sessionId, err = srv.dbo.CreateSession(userId, p.Auth.UserIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.DeletePreSessionByRequestId(preSession.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var tokenString simple.WebTokenString
	tokenString, err = srv.jwtkm.MakeJWToken(userId, sessionId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_JWTCreation, fmt.Sprintf(RpcErrorMsgF_JWTCreation, err.Error()), nil)
	}

	// Journaling.
	logEvent := le.NewLogEventWithFields(let.NewLogEventTypeWithValue(ev.NewEnumValue(let.LogEventType_LogIn)), userId, p.Email, p.Auth.UserIPAB, nil)

	err = srv.dbo.SaveLogEvent(logEvent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.LogUserInResult{NextStep: 0, IsWebTokenSet: true, WebTokenString: tokenString}

	return result, nil
}

func (srv *Server) logUserOut(p *rpc2.LogUserOutParams) (result *rpc2.LogUserOutResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	err := srv.dbo.DeleteSession(thisUserData.GetSession().Id, thisUserData.GetUser().GetUserParameters().GetId(), thisUserData.GetSession().UserIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Journaling.
	logEvent := le.NewLogEventWithFields(let.NewLogEventTypeWithValue(ev.NewEnumValue(let.LogEventType_LogOut)), thisUserData.GetUser().GetUserParameters().GetId(), thisUserData.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB, nil)

	err = srv.dbo.SaveLogEvent(logEvent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.LogUserOutResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

func (srv *Server) logUserOutA(p *rpc2.LogUserOutAParams) (result *rpc2.LogUserOutAResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var callerData derived1.IUserData
	callerData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !callerData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var user derived1.IUser
	var err error
	user, err = srv.dbo.GetUserById(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if user == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIsNotFound, RpcErrorMsg_UserIsNotFound, p.UserId)
	}

	if user.GetUserParameters().GetId() != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_DatabaseInconsistency, RpcErrorMsg_DatabaseInconsistency, user)
	}

	var session *am.Session
	session, err = srv.dbo.GetSessionByUserId(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if session == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SessionIsNotFound, RpcErrorMsg_SessionIsNotFound, p.UserId)
	}

	if session.UserId != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_DatabaseInconsistency, RpcErrorMsg_DatabaseInconsistency, session)
	}

	err = srv.dbo.DeleteSessionByUserId(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Journaling.
	logEvent := le.NewLogEventWithFields(let.NewLogEventTypeWithValue(ev.NewEnumValue(let.LogEventType_LogOutA)), session.UserId, user.GetUserParameters().GetEmail(), session.UserIPAB, callerData.GetUser().GetUserParameters().GetIdPtr())

	err = srv.dbo.SaveLogEvent(logEvent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.LogUserOutAResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

func (srv *Server) getListOfLoggedUsers(p *rpc2.GetListOfLoggedUsersParams) (result *rpc2.GetListOfLoggedUsersResult, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	_, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	result = &rpc2.GetListOfLoggedUsersResult{}

	var err error
	result.LoggedUserIds, err = srv.dbo.GetListOfLoggedUsers()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return result, nil
}

func (srv *Server) getListOfLoggedUsersOnPage(p *rpc2.GetListOfLoggedUsersOnPageParams) (result *rpc2.GetListOfLoggedUsersOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	_, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	var loggedUsersCount base2.Count
	var err error
	loggedUsersCount, err = srv.dbo.CountLoggedUsers()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var userIdsOnPage []base2.Id
	userIdsOnPage, err = srv.dbo.GetListOfLoggedUsersOnPage(p.Page, srv.settings.SystemSettings.PageSize)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.GetListOfLoggedUsersOnPageResult{
		LoggedUserIds: userIdsOnPage,
		PageData: &rpc3.PageData{
			PageNumber:  p.Page,
			TotalPages:  base2.CalculateTotalPages(loggedUsersCount, srv.settings.SystemSettings.PageSize),
			PageSize:    srv.settings.SystemSettings.PageSize,
			ItemsOnPage: base2.Count(len(userIdsOnPage)),
			TotalItems:  loggedUsersCount,
		},
	}

	return result, nil
}

func (srv *Server) getListOfAllUsers(p *rpc2.GetListOfAllUsersParams) (result *rpc2.GetListOfAllUsersResult, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var userIds []base2.Id
	var err error
	userIds, err = srv.dbo.GetListOfAllUsers()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.GetListOfAllUsersResult{
		UserIds: userIds,
	}

	return result, nil
}

func (srv *Server) getListOfAllUsersOnPage(p *rpc2.GetListOfAllUsersOnPageParams) (result *rpc2.GetListOfAllUsersOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	allUsersCount, err := srv.dbo.CountAllUsers()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var userIdsOnPage []base2.Id
	userIdsOnPage, err = srv.dbo.GetListOfAllUsersOnPage(p.Page, srv.settings.SystemSettings.PageSize)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.GetListOfAllUsersOnPageResult{
		UserIds: userIdsOnPage,
		PageData: &rpc3.PageData{
			PageNumber:  p.Page,
			TotalPages:  base2.CalculateTotalPages(allUsersCount, srv.settings.SystemSettings.PageSize),
			PageSize:    srv.settings.SystemSettings.PageSize,
			ItemsOnPage: base2.Count(len(userIdsOnPage)),
			TotalItems:  allUsersCount,
		},
	}

	return result, nil
}

func (srv *Server) isUserLoggedIn(p *rpc2.IsUserLoggedInParams) (result *rpc2.IsUserLoggedInResult, re *jrm1.RpcError) {
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	_, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	result = &rpc2.IsUserLoggedInResult{
		UserId: p.UserId,
	}

	sessionsCount, err := srv.dbo.CountSessionsByUserId(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if sessionsCount == 1 {
		result.IsUserLoggedIn = true
	} else {
		result.IsUserLoggedIn = false
	}

	return result, nil
}

// Various actions.

func (srv *Server) changePassword(p *rpc2.ChangePasswordParams) (result *rpc2.ChangePasswordResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	switch p.StepN {
	case 1:
		return srv.changePasswordStep1(p, thisUserData)
	case 2:
		return srv.changePasswordStep2(p, thisUserData)
	default:
		// Step is not supported.
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_StepIsUnknown, RpcErrorMsg_StepIsUnknown, nil)
	}
}

func (srv *Server) changePasswordStep1(p *rpc2.ChangePasswordParams, ud derived1.IUserData) (result *rpc2.ChangePasswordResult, re *jrm1.RpcError) {
	// Is new password set ?
	if len(p.NewPassword) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NewPasswordIsNotSet, RpcErrorMsg_NewPasswordIsNotSet, nil)
	}

	var pc = &am.PasswordChange{
		UserId:   ud.GetUser().GetUserParameters().GetId(),
		UserIPAB: p.Auth.UserIPAB,
	}

	// Check for duplicate request.
	passwordChangesCount, err := srv.dbo.CountPasswordChangesByUserId(ud.GetUser().GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if passwordChangesCount > 0 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PasswordChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserAlreadyStartedToChangePassword, RpcErrorMsg_UserAlreadyStartedToChangePassword, nil)
	}

	// Check the new password.
	re = isPasswordAllowed(p.NewPassword)
	if re != nil {
		return nil, re
	}

	pc.NewPasswordBytes, err = bpp.PackSymbols([]rune(p.NewPassword))
	if err != nil {
		// This error is very unlikely to happen.
		// So if it occurs, then it is an anomaly.
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_BPP_PackSymbols, fmt.Sprintf(RpcErrorMsgF_BPP_PackSymbols, err.Error()), nil)
	}

	if len(pc.NewPasswordBytes) > int(srv.settings.SystemSettings.UserPasswordMaxLenInBytes) {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordIsTooLong, RpcErrorMsg_PasswordIsTooLong, nil)
	}

	// Request ID.
	pc.RequestId, re = srv.createRequestIdForPasswordChange()
	if re != nil {
		return nil, re
	}

	// Password salt.
	pc.AuthDataBytes, err = bpp.GenerateRandomSalt()
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_BPP_GenerateRandomSalt, RpcErrorMsg_BPP_GenerateRandomSalt, nil)
	}

	// Verification code.
	pc.VerificationCode, re = srv.createVerificationCode()
	if re != nil {
		return nil, re
	}

	re = srv.sendVerificationCodeForPwdChange(ud.GetUser().GetUserParameters().GetEmail(), *pc.VerificationCode)
	if re != nil {
		return nil, re
	}

	// Captcha.
	pc.IsCaptchaRequired, err = srv.isCaptchaNeededForPasswordChange(ud.GetUser().GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var captchaData *rm.CreateCaptchaResult
	if pc.IsCaptchaRequired {
		captchaData, re = srv.createCaptcha()
		if re != nil {
			return nil, re
		}

		pc.CaptchaId = (*simple.CaptchaId)(&captchaData.TaskId)
	}

	// Save the request.
	err = srv.dbo.CreatePasswordChangeRequest(pc)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Response.
	result = &rpc2.ChangePasswordResult{
		NextStep:      2,
		RequestId:     *pc.RequestId,
		AuthDataBytes: pc.AuthDataBytes,
	}

	if pc.IsCaptchaRequired {
		result.IsCaptchaNeeded = true
		result.CaptchaId = *pc.CaptchaId
	} else {
		result.IsCaptchaNeeded = false
	}

	return result, nil
}

func (srv *Server) changePasswordStep2(p *rpc2.ChangePasswordParams, ud derived1.IUserData) (result *rpc2.ChangePasswordResult, re *jrm1.RpcError) {
	// Check input parameters.
	if len(p.RequestId) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RequestIdIsNotSet, RpcErrorMsg_RequestIdIsNotSet, nil)
	}
	if len(p.AuthChallengeResponse) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_AuthChallengeResponseIsNotSet, RpcErrorMsg_AuthChallengeResponseIsNotSet, nil)
	}
	if len(p.VerificationCode) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsNotSet, RpcErrorMsg_VerificationCodeIsNotSet, nil)
	}

	// Check for duplicate request.
	passwordChangesCount, err := srv.dbo.CountPasswordChangesByUserId(ud.GetUser().GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if passwordChangesCount != 1 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PasswordChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserAlreadyStartedToChangePassword, RpcErrorMsg_UserAlreadyStartedToChangePassword, nil)
	}

	// Get the request for password change.
	var pcr *am.PasswordChange
	pcr, err = srv.dbo.GetPasswordChangeByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if pcr == nil {
		// Request Id code can not be guessed.
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PasswordChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordChangeIsNotFound, RpcErrorMsg_PasswordChangeIsNotFound, nil)
	}

	// Check the Request ID (indirectly).
	if (pcr.UserId != ud.GetUser().GetUserParameters().GetId()) || (pcr.RequestId == nil) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PasswordChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordChangeIsNotFound, RpcErrorMsg_PasswordChangeIsNotFound, nil)
	}

	// Check the IP address.
	if !p.Auth.UserIPAB.Equal(pcr.UserIPAB) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PasswordChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordChangeIsNotFound, RpcErrorMsg_PasswordChangeIsNotFound, nil)
	}

	// Check the captcha answer if needed.
	var ccr *rm.CheckCaptchaResult
	if pcr.IsCaptchaRequired {
		if len(p.CaptchaAnswer) == 0 {
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CaptchaAnswerIsNotSet, RpcErrorMsg_CaptchaAnswerIsNotSet, nil)
		}
		if pcr.CaptchaId == nil {
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CaptchaIdIsNotSet, RpcErrorMsg_CaptchaIdIsNotSet, nil)
		}

		// Check captcha answer.
		ccr, re = srv.checkCaptcha(*pcr.CaptchaId, p.CaptchaAnswer)
		if re != nil {
			return nil, re
		}

		if !ccr.IsSuccess {
			srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_CaptchaAnswerMismatch), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

			err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
			if err != nil {
				return nil, srv.databaseError(err)
			}

			// When captcha guess is wrong, we delete the password change
			// request to start the process from the first step.
			err = srv.dbo.DeletePasswordChangeByRequestId(*pcr.RequestId)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CaptchaAnswerIsWrong, RpcErrorMsg_CaptchaAnswerIsWrong, nil)
		}
	}

	// Check the password.
	var ok bool
	ok, re = srv.checkPassword(ud.GetUser().GetUserParameters().GetId(), pcr.AuthDataBytes, p.AuthChallengeResponse)
	if re != nil {
		return nil, re
	}

	if !ok {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PasswordMismatch), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// When password is wrong, we delete the password change request to
		// start the process from the first step.
		err = srv.dbo.DeletePasswordChangeByRequestId(*pcr.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordIsWrong, RpcErrorMsg_PasswordIsWrong, nil)
	}

	// Check the verification code.
	ok, err = srv.dbo.CheckVerificationCodeForPwdChange(p.RequestId, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		// Verification code can not be guessed.
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_VerificationCodeMismatch), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// Delete the password change request on error to avoid brute force
		// checks.
		err = srv.dbo.DeletePasswordChangeByRequestId(*pcr.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsWrong, RpcErrorMsg_VerificationCodeIsWrong, nil)
	}

	// Update request's flags.
	pcvf := &am.PasswordChangeVerificationFlags{
		IsVerifiedByPassword: true,
		IsVerifiedByEmail:    true,
	}
	if pcr.IsCaptchaRequired {
		x := true
		pcvf.IsVerifiedByCaptcha = (*base2.Flag)(&x)
	}

	err = srv.dbo.SetPasswordChangeVFlags(ud.GetUser().GetUserParameters().GetId(), *pcr.RequestId, pcvf)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Change the password.
	err = srv.dbo.SetUserPassword(ud.GetUser().GetUserParameters().GetId(), ud.GetUser().GetUserParameters().GetEmail(), pcr.NewPasswordBytes)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Close the session to force re-logging.
	err = srv.dbo.DeleteSession(ud.GetSession().Id, ud.GetUser().GetUserParameters().GetId(), p.Auth.UserIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// N.B.: We do not immediately delete the request for password change here
	// to avoid spamming with these requests. It will be deleted later by the
	// scheduler.

	// Response.
	result = &rpc2.ChangePasswordResult{
		Success: rpc3.Success{
			OK: true,
		},
		NextStep: 0,
	}
	return result, nil
}

func (srv *Server) changeEmail(p *rpc2.ChangeEmailParams) (result *rpc2.ChangeEmailResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	switch p.StepN {
	case 1:
		return srv.changeEmailStep1(p, thisUserData)
	case 2:
		return srv.changeEmailStep2(p, thisUserData)
	default:
		// Step is not supported.
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_StepIsUnknown, RpcErrorMsg_StepIsUnknown, nil)
	}
}

func (srv *Server) changeEmailStep1(p *rpc2.ChangeEmailParams, ud derived1.IUserData) (result *rpc2.ChangeEmailResult, re *jrm1.RpcError) {
	// Is new e-mail address set ?
	if len(p.NewEmail) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NewEmailIsNotSet, RpcErrorMsg_NewEmailIsNotSet, nil)
	}

	// Check for duplicate request.
	emailChangesCount, err := srv.dbo.CountEmailChangesByUserId(ud.GetUser().GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if emailChangesCount > 0 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_EmailChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserAlreadyStartedToChangeEmail, RpcErrorMsg_UserAlreadyStartedToChangeEmail, nil)
	}

	var ec = &am.EmailChange{
		UserId:   ud.GetUser().GetUserParameters().GetId(),
		UserIPAB: p.Auth.UserIPAB,
		NewEmail: p.NewEmail,
	}

	// Check the new e-mail address.
	// Is e-mail address valid ?
	re = srv.isEmailOfUserValid(ec.NewEmail)
	if re != nil {
		return nil, re
	}

	// Check that the new e-mail address is not used.
	var usersCount base2.Count
	usersCount, err = srv.dbo.CountUsersWithEmail(ec.NewEmail)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	isAddrFree := usersCount == 0
	if !isAddrFree {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsUsed, RpcErrorMsg_EmailAddressIsUsed, nil)
	}

	// Request ID.
	ec.RequestId, re = srv.createRequestIdForPasswordChange()
	if re != nil {
		return nil, re
	}

	// Password salt.
	ec.AuthDataBytes, err = bpp.GenerateRandomSalt()
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_BPP_GenerateRandomSalt, RpcErrorMsg_BPP_GenerateRandomSalt, nil)
	}

	// Verification code for the old e-mail address.
	ec.VerificationCodeOld, re = srv.createVerificationCode()
	if re != nil {
		return nil, re
	}

	re = srv.sendVerificationCodeForEmailChange(ud.GetUser().GetUserParameters().GetEmail(), *ec.VerificationCodeOld)
	if re != nil {
		return nil, re
	}

	// Verification code for the new e-mail address.
	ec.VerificationCodeNew, re = srv.createVerificationCode()
	if re != nil {
		return nil, re
	}

	re = srv.sendVerificationCodeForEmailChange(ec.NewEmail, *ec.VerificationCodeNew)
	if re != nil {
		return nil, re
	}

	// Captcha.
	ec.IsCaptchaRequired, err = srv.isCaptchaNeededForEmailChange(ud.GetUser().GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var captchaData *rm.CreateCaptchaResult
	if ec.IsCaptchaRequired {
		captchaData, re = srv.createCaptcha()
		if re != nil {
			return nil, re
		}

		ec.CaptchaId = (*simple.CaptchaId)(&captchaData.TaskId)
	}

	// Save the request.
	err = srv.dbo.CreateEmailChangeRequest(ec)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Response.
	result = &rpc2.ChangeEmailResult{
		NextStep:      2,
		RequestId:     *ec.RequestId,
		AuthDataBytes: ec.AuthDataBytes,
	}

	if ec.IsCaptchaRequired {
		result.IsCaptchaNeeded = true
		result.CaptchaId = *ec.CaptchaId
	} else {
		result.IsCaptchaNeeded = false
	}

	return result, nil
}

func (srv *Server) changeEmailStep2(p *rpc2.ChangeEmailParams, ud derived1.IUserData) (result *rpc2.ChangeEmailResult, re *jrm1.RpcError) {
	// Check input parameters.
	if len(p.RequestId) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RequestIdIsNotSet, RpcErrorMsg_RequestIdIsNotSet, nil)
	}
	if len(p.AuthChallengeResponse) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_AuthChallengeResponseIsNotSet, RpcErrorMsg_AuthChallengeResponseIsNotSet, nil)
	}
	if len(p.VerificationCodeOld) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsNotSet, RpcErrorMsg_VerificationCodeIsNotSet, nil)
	}
	if len(p.VerificationCodeNew) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsNotSet, RpcErrorMsg_VerificationCodeIsNotSet, nil)
	}

	// Check for duplicate request.
	emailChangesCount, err := srv.dbo.CountEmailChangesByUserId(ud.GetUser().GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if emailChangesCount != 1 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_EmailChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserAlreadyStartedToChangeEmail, RpcErrorMsg_UserAlreadyStartedToChangeEmail, nil)
	}

	// Get the request for e-mail address change.
	var ecr *am.EmailChange
	ecr, err = srv.dbo.GetEmailChangeByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if ecr == nil {
		// Request Id code can not be guessed.
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_EmailChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailChangeIsNotFound, RpcErrorMsg_EmailChangeIsNotFound, nil)
	}

	// Check the Request ID (indirectly).
	if (ecr.UserId != ud.GetUser().GetUserParameters().GetId()) || (ecr.RequestId == nil) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_EmailChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailChangeIsNotFound, RpcErrorMsg_EmailChangeIsNotFound, nil)
	}

	// Check the IP address.
	if !p.Auth.UserIPAB.Equal(ecr.UserIPAB) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_EmailChangeHacking), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_EmailChangeIsNotFound, RpcErrorMsg_EmailChangeIsNotFound, nil)
	}

	// Check the captcha answer if needed.
	var ccr *rm.CheckCaptchaResult
	if ecr.IsCaptchaRequired {
		if len(p.CaptchaAnswer) == 0 {
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CaptchaAnswerIsNotSet, RpcErrorMsg_CaptchaAnswerIsNotSet, nil)
		}
		if ecr.CaptchaId == nil {
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CaptchaIdIsNotSet, RpcErrorMsg_CaptchaIdIsNotSet, nil)
		}

		// Check captcha answer.
		ccr, re = srv.checkCaptcha(*ecr.CaptchaId, p.CaptchaAnswer)
		if re != nil {
			return nil, re
		}

		if !ccr.IsSuccess {
			srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_CaptchaAnswerMismatch), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

			err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
			if err != nil {
				return nil, srv.databaseError(err)
			}

			// When captcha guess is wrong, we delete the e-mail address change
			// request to start the process from the first step.
			err = srv.dbo.DeleteEmailChangeByRequestId(*ecr.RequestId)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_CaptchaAnswerIsWrong, RpcErrorMsg_CaptchaAnswerIsWrong, nil)
		}
	}

	// Check the password.
	var ok bool
	ok, re = srv.checkPassword(ud.GetUser().GetUserParameters().GetId(), ecr.AuthDataBytes, p.AuthChallengeResponse)
	if re != nil {
		return nil, re
	}

	if !ok {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_PasswordMismatch), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// When password is wrong, we delete the password change request to
		// start the process from the first step.
		err = srv.dbo.DeleteEmailChangeByRequestId(*ecr.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordIsWrong, RpcErrorMsg_PasswordIsWrong, nil)
	}

	// Check the verification codes.
	ok, err = srv.dbo.CheckVerificationCodesForEmailChange(p.RequestId, p.VerificationCodeOld, p.VerificationCodeNew)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		// Verification codes can not be guessed.
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_VerificationCodeMismatch), ud.GetUser().GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.GetUser().GetUserParameters().GetId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// Delete the e-mail address change request on error to avoid brute
		// force checks.
		err = srv.dbo.DeleteEmailChangeByRequestId(*ecr.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeIsWrong, RpcErrorMsg_VerificationCodeIsWrong, nil)
	}

	// Update request's flags.
	ecvf := &am.EmailChangeVerificationFlags{
		IsVerifiedByPassword: true,
		IsVerifiedByOldEmail: true,
		IsVerifiedByNewEmail: true,
	}
	if ecr.IsCaptchaRequired {
		x := true
		ecvf.IsVerifiedByCaptcha = (*base2.Flag)(&x)
	}

	err = srv.dbo.SetEmailChangeVFlags(ud.GetUser().GetUserParameters().GetId(), *ecr.RequestId, ecvf)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Change the e-mail address.
	err = srv.dbo.SetUserEmail(ud.GetUser().GetUserParameters().GetId(), ud.GetUser().GetUserParameters().GetEmail(), ecr.NewEmail)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Close the session to force re-logging.
	err = srv.dbo.DeleteSession(ud.GetSession().Id, ud.GetUser().GetUserParameters().GetId(), p.Auth.UserIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// N.B.: We do not immediately delete the request for e-mail address change
	// here to avoid spamming with these requests. It will be deleted later by
	// the scheduler.

	// Response.
	result = &rpc2.ChangeEmailResult{
		Success: rpc3.Success{
			OK: true,
		},
		NextStep: 0,
	}
	return result, nil
}

func (srv *Server) getUserSession(p *rpc2.GetUserSessionParams) (result *rpc2.GetUserSessionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var callerData derived1.IUserData
	callerData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !callerData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var err error
	var session *am.Session
	session, err = srv.dbo.GetSessionByUserId(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if session == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SessionIsNotFound, RpcErrorMsg_SessionIsNotFound, nil)
	}

	user := u.NewUser()
	userParams := user.GetUserParameters()
	userParams.SetId(p.UserId)
	user.SetUserParameters(userParams)

	result = &rpc2.GetUserSessionResult{
		User:    user,
		Session: session,
	}

	return result, nil
}

// User properties.

func (srv *Server) getUserName(p *rpc2.GetUserNameParams) (result *rpc2.GetUserNameResult, re *jrm1.RpcError) {
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	_, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	var err error
	var userName *simple.Name
	userName, err = srv.dbo.GetUserNameById(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if userName == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserNameIsNotFound, RpcErrorMsg_UserNameIsNotFound, p.UserId)
	}

	user := u.NewUser()
	userParams := user.GetUserParameters()
	userParams.SetId(p.UserId)
	userParams.SetName(*userName)
	user.SetUserParameters(userParams)

	result = &rpc2.GetUserNameResult{
		User: user,
	}

	return result, nil
}

func (srv *Server) getUserRoles(p *rpc2.GetUserRolesParams) (result *rpc2.GetUserRolesResult, re *jrm1.RpcError) {
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	_, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	var err error
	var roles *ur.UserRoles
	roles, err = srv.dbo.GetUserRolesById(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if roles == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIsNotFound, RpcErrorMsg_UserIsNotFound, nil)
	}

	roles.IsAdministrator = srv.isUserAdministrator(p.UserId)
	roles.IsModerator = srv.isUserModerator(p.UserId)

	user := u.NewUser()
	userParams := user.GetUserParameters()
	userParams.SetId(p.UserId)
	userParams.SetRoles(roles)
	user.SetUserParameters(userParams)

	result = &rpc2.GetUserRolesResult{
		User: user,
	}

	return result, nil
}

func (srv *Server) viewUserParameters(p *rpc2.ViewUserParametersParams) (result *rpc2.ViewUserParametersResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var err error
	var userParameters base22.IUserParameters
	userParameters, err = srv.dbo.ViewUserParametersById(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if userParameters == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIsNotFound, RpcErrorMsg_UserIsNotFound, nil)
	}

	roles := userParameters.GetRoles()
	roles.IsModerator = srv.isUserModerator(p.UserId)
	roles.IsAdministrator = srv.isUserAdministrator(p.UserId)
	userParameters.SetRoles(roles)

	user := u.NewUser()
	user.SetUserParameters(userParameters)

	result = &rpc2.ViewUserParametersResult{
		User: user,
	}

	return result, nil
}

func (srv *Server) setUserRoleAuthor(p *rpc2.SetUserRoleAuthorParams) (result *rpc2.SetUserRoleAuthorResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	err := srv.dbo.SetUserRoleAuthor(p.UserId, p.IsRoleEnabled)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.SetUserRoleAuthorResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

func (srv *Server) setUserRoleWriter(p *rpc2.SetUserRoleWriterParams) (result *rpc2.SetUserRoleWriterResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	err := srv.dbo.SetUserRoleWriter(p.UserId, p.IsRoleEnabled)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.SetUserRoleWriterResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

func (srv *Server) setUserRoleReader(p *rpc2.SetUserRoleReaderParams) (result *rpc2.SetUserRoleReaderResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsAdministrator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	err := srv.dbo.SetUserRoleReader(p.UserId, p.IsRoleEnabled)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.SetUserRoleReaderResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

func (srv *Server) getSelfRoles(p *rpc2.GetSelfRolesParams) (result *rpc2.GetSelfRolesResult, re *jrm1.RpcError) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	user := u.NewUser()
	user.SetUserParameters(thisUserData.GetUser().GetUserParameters())

	result = &rpc2.GetSelfRolesResult{
		User: user,
	}

	return result, nil
}

// User banning.

func (srv *Server) banUser(p *rpc2.BanUserParams) (result *rpc2.BanUserResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsModerator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	err := srv.dbo.SetUserRoleCanLogIn(p.UserId, false)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.UpdateUserBanTime(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var sessionsCount base2.Count
	sessionsCount, err = srv.dbo.CountSessionsByUserId(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if sessionsCount > 0 {
		err = srv.dbo.DeleteSessionByUserId(p.UserId)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &rpc2.BanUserResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

func (srv *Server) unbanUser(p *rpc2.UnbanUserParams) (result *rpc2.UnbanUserResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData derived1.IUserData
	thisUserData, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !thisUserData.GetUser().GetUserParameters().GetRoles().IsModerator {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", p.Auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	err := srv.dbo.SetUserRoleCanLogIn(p.UserId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.UnbanUserResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// Other.

func (srv *Server) showDiagnosticData() (result *rpc2.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	trc, src := srv.js.GetRequestsCount()

	result = &rpc2.ShowDiagnosticDataResult{
		RequestsCount: rpc3.RequestsCount{
			TotalRequestsCount:      base2.Text(trc),
			SuccessfulRequestsCount: base2.Text(src),
		},
	}

	return result, nil
}

func (srv *Server) test(_ *rpc2.TestParams) (result *rpc2.TestResult, re *jrm1.RpcError) {
	result = &rpc2.TestResult{}

	time.Sleep(time.Second * 10)

	return result, nil
}
