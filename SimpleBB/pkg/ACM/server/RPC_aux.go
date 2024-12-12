package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/rpc"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/IncidentType"
	am "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/UserData"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/models/net"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	server2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"log"
	"net/mail"
	"time"

	bpp "github.com/vault-thirteen/BytePackedPassword"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	rc "github.com/vault-thirteen/SimpleBB/pkg/RCS/client"
	sc "github.com/vault-thirteen/SimpleBB/pkg/SMTP/client"
	num "github.com/vault-thirteen/auxie/number"
)

// Auxiliary functions used in RPC functions.

const (
	ErrAuthData = "authorisation data error"
)

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if err == nil {
		return
	}

	if srv.settings.SystemSettings.IsDebugMode {
		log.Println(err)
	}
}

// processDatabaseError processes a database error.
func (srv *Server) processDatabaseError(err error) {
	if err == nil {
		return
	}

	if server2.IsNetworkError(err) {
		log.Println(fmt.Sprintf(server2.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return
}

// databaseError processes the database error and returns an RPC error.
func (srv *Server) databaseError(err error) (re *jrm1.RpcError) {
	srv.processDatabaseError(err)
	return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Database, server2.RpcErrorMsg_Database, err)
}

// Token-related functions.

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *cmr.Auth) (re *jrm1.RpcError) {
	if (auth == nil) ||
		(len(auth.UserIPA) == 0) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", nil)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Authorisation, server2.RpcErrorMsg_Authorisation, nil)
	}

	var err error
	auth.UserIPAB, err = cn.ParseIPA(auth.UserIPA)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Authorisation, server2.RpcErrorMsg_Authorisation, nil)
	}

	return nil
}

// mustBeNoAuthToken ensures that an authorisation token is not present. If the
// token is present, an error is returned and the caller of this function must
// stop and return this error.
func (srv *Server) mustBeNoAuthToken(auth *cmr.Auth) (re *jrm1.RpcError) {
	re = srv.mustBeAuthUserIPA(auth)
	if re != nil {
		return re
	}

	if len(auth.Token) > 0 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", auth.UserIPAB)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Permission, server2.RpcErrorMsg_Permission, nil)
	}

	return nil
}

// mustBeAnAuthToken ensures that an authorisation token is present and is
// valid. If the token is absent or invalid, an error is returned and the caller
// of this function must stop and return this error. User data is returned when
// token is valid.
func (srv *Server) mustBeAnAuthToken(auth *cmr.Auth) (ud derived1.IUserData, re *jrm1.RpcError) {
	re = srv.mustBeAuthUserIPA(auth)
	if re != nil {
		return nil, re
	}

	if len(auth.Token) == 0 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_IllegalAccessAttempt), "", auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Authorisation, server2.RpcErrorMsg_Authorisation, nil)
	}

	var err error
	ud, err = srv.getUserDataByAuthToken(auth.Token)
	if err != nil {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_FakeToken), "", auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Authorisation, server2.RpcErrorMsg_Authorisation, nil)
	}

	if bytes.Compare(auth.UserIPAB, ud.GetSession().UserIPAB) != 0 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(cm.IncidentType_FakeIPA), "", auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Authorisation, server2.RpcErrorMsg_Authorisation, nil)
	}

	return ud, nil
}

// getUserDataByAuthToken validates the token and returns information about the
// user and its session. When 'userData' is set (not null), all its fields are
// also set (not null).
func (srv *Server) getUserDataByAuthToken(authToken simple.WebTokenString) (userData derived1.IUserData, err error) {
	var userId, sessionId base2.Id
	userId, sessionId, err = srv.jwtkm.ValidateToken(authToken)
	if err != nil {
		return nil, err
	}

	userData = am.NewUserData()

	var x derived1.IUser
	x, err = srv.dbo.GetUserById(userId)
	if err != nil {
		srv.processDatabaseError(err)
		return nil, err
	}

	userData.SetUser(x)

	if (userData.GetUser() == nil) ||
		(userData.GetUser().GetUserParameters().GetId() != userId) ||
		(!userData.GetUser().GetUserParameters().GetRoles().CanLogIn) {
		return nil, errors.New(ErrAuthData)
	}

	// Attach special user roles from settings.
	roles := userData.GetUser().GetUserParameters().GetRoles()
	roles.IsModerator = srv.isUserModerator(userData.GetUser().GetUserParameters().GetId())
	roles.IsAdministrator = srv.isUserAdministrator(userData.GetUser().GetUserParameters().GetId())

	uParams := userData.GetUser().GetUserParameters()
	uParams.SetRoles(roles)
	userData.GetUser().SetUserParameters(uParams)

	var session *models.Session
	session, err = srv.dbo.GetSessionByUserId(userId)
	if err != nil {
		srv.processDatabaseError(err)
		return nil, err
	}

	userData.SetSession(session)

	if (userData.GetSession() == nil) ||
		(userData.GetSession().Id != sessionId) {
		return nil, errors.New(ErrAuthData)
	}

	return userData, nil
}

// Other functions.

func (srv *Server) checkCaptcha(captchaId simple.CaptchaId, answer simple.CaptchaAnswer) (result *rm.CheckCaptchaResult, re *jrm1.RpcError) {
	var params = rm.CheckCaptchaParams{TaskId: captchaId.ToString()}
	var err error

	if len(answer) == 0 {
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Captcha, server2.RpcErrorMsg_Captcha, nil)
	}

	params.Value, err = num.ParseUint(answer.ToString())
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Captcha, server2.RpcErrorMsg_Captcha, nil)
	}

	result = new(rm.CheckCaptchaResult)
	re, err = srv.rcsServiceClient.MakeRequest(context.Background(), rc.FuncCheckCaptcha, params, result)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Captcha, server2.RpcErrorMsg_Captcha, nil)
	}

	return result, nil
}

func (srv *Server) checkPassword(userId base2.Id, salt []byte, userKey []byte) (ok bool, re *jrm1.RpcError) {
	var pwd *[]byte
	var err error
	pwd, err = srv.dbo.GetUserPasswordById(userId)
	if err != nil {
		return false, srv.databaseError(err)
	}
	if pwd == nil {
		return false, nil
	}

	var pwdRunes []rune
	pwdRunes, err = bpp.UnpackBytes(*pwd)
	if err != nil {
		srv.logError(err)
		return false, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Password, server2.RpcErrorMsg_Password, nil)
	}

	ok, err = bpp.CheckHashKey(string(pwdRunes), salt, userKey)
	if err != nil {
		srv.logError(err)
		return false, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Password, server2.RpcErrorMsg_Password, nil)
	}

	return ok, nil
}

func (srv *Server) createCaptcha() (result *rm.CreateCaptchaResult, re *jrm1.RpcError) {
	var params = rm.CreateCaptchaParams{}

	result = new(rm.CreateCaptchaResult)
	var err error
	re, err = srv.rcsServiceClient.MakeRequest(context.Background(), rc.FuncCreateCaptcha, params, result)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Captcha, server2.RpcErrorMsg_Captcha, nil)
	}
	if result.IsImageDataReturned {
		// We do not return an image in a JSON message.
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Captcha, server2.RpcErrorMsg_Captcha, nil)
	}

	return result, nil
}

func (srv *Server) createRequestIdForLogIn() (rid *simple.RequestId, re *jrm1.RpcError) {
	s, err := srv.ridg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RequestIdGenerator, RpcErrorMsg_RequestIdGenerator, nil)
	}

	return (*simple.RequestId)(s), nil
}

func (srv *Server) createRequestIdForPasswordChange() (rid *simple.RequestId, re *jrm1.RpcError) {
	s, err := srv.ridg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RequestIdGenerator, RpcErrorMsg_RequestIdGenerator, nil)
	}

	return (*simple.RequestId)(s), nil
}

func (srv *Server) createVerificationCode() (vc *simple.VerificationCode, re *jrm1.RpcError) {
	var err error
	var s *string
	s, err = srv.vcg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeGenerator, RpcErrorMsg_VerificationCodeGenerator, nil)
	}

	return (*simple.VerificationCode)(s), nil
}

func (srv *Server) isCaptchaNeededForLogIn(email simple.Email) (isCaptchaNeeded base2.Flag, err error) {
	var lastBadLogInTime *time.Time
	lastBadLogInTime, err = srv.dbo.GetUserLastBadLogInTimeByEmail(email)
	if err != nil {
		return false, err
	}
	if lastBadLogInTime == nil {
		return false, nil
	}

	if time.Now().Before(lastBadLogInTime.Add(time.Duration(srv.settings.SystemSettings.LogInTryTimeout) * time.Second)) {
		return true, nil
	}

	return false, nil
}

func (srv *Server) isCaptchaNeededForEmailChange(userId base2.Id) (isCaptchaNeeded base2.Flag, err error) {
	var lastBadActionTime *time.Time
	lastBadActionTime, err = srv.dbo.GetUserLastBadActionTimeById(userId)
	if err != nil {
		return false, err
	}
	if lastBadActionTime == nil {
		return false, nil
	}

	if time.Now().Before(lastBadActionTime.Add(time.Duration(srv.settings.SystemSettings.ActionTryTimeout) * time.Second)) {
		return true, nil
	}

	return false, nil
}

func (srv *Server) isCaptchaNeededForPasswordChange(userId base2.Id) (isCaptchaNeeded base2.Flag, err error) {
	var lastBadActionTime *time.Time
	lastBadActionTime, err = srv.dbo.GetUserLastBadActionTimeById(userId)
	if err != nil {
		return false, err
	}
	if lastBadActionTime == nil {
		return false, nil
	}

	if time.Now().Before(lastBadActionTime.Add(time.Duration(srv.settings.SystemSettings.ActionTryTimeout) * time.Second)) {
		return true, nil
	}

	return false, nil
}

func (srv *Server) isEmailOfUserValid(email simple.Email) (re *jrm1.RpcError) {
	_, err := mail.ParseAddress(email.ToString())
	if err != nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsNotValid, RpcErrorMsg_EmailAddressIsNotValid, nil)

	}

	return nil
}

func isPasswordAllowed(password simple.Password) (re *jrm1.RpcError) {
	_, err := bpp.IsPasswordAllowed(password.ToString())
	if err != nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordIsNotValid, RpcErrorMsg_PasswordIsNotValid, nil)
	}

	return nil
}

func (srv *Server) isUserAdministrator(userId base2.Id) (isAdministrator base2.Flag) {
	// While system has only few administrators, the simple array look-up is
	// faster than access to a map.

	for _, id := range srv.settings.UserRoleSettings.AdministratorIds {
		if id == userId {
			return true
		}
	}

	return false
}

func (srv *Server) isUserModerator(userId base2.Id) (isModerator base2.Flag) {
	// While system has only few moderators, the simple array look-up is
	// faster than access to a map.

	for _, id := range srv.settings.UserRoleSettings.ModeratorIds {
		if id == userId {
			return true
		}
	}

	return false
}

func (srv *Server) sendEmailMessage(params sm.SendMessageParams) (re *jrm1.RpcError) {
	var result = new(sm.SendMessageResult)

	var err error
	re, err = srv.smtpServiceClient.MakeRequest(context.Background(), sc.FuncSendMessage, params, result)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_SmtpModule, RpcErrorMsg_SmtpModule, nil)
	}

	return nil
}

func (srv *Server) sendVerificationCodeForReg(email simple.Email, code simple.VerificationCode) (re *jrm1.RpcError) {
	var subject = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplateForRegVCode.ToString(), srv.settings.SystemSettings.SiteName))
	var msg = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForRegVCode.ToString(), srv.settings.SystemSettings.SiteName, code))
	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}
	return srv.sendEmailMessage(params)
}

func (srv *Server) sendVerificationCodeForLogIn(email simple.Email, code simple.VerificationCode) (re *jrm1.RpcError) {
	var subject = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplateForRegVCode.ToString(), srv.settings.SystemSettings.SiteName))
	var msg = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForLogIn.ToString(), code))
	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}
	return srv.sendEmailMessage(params)
}

func (srv *Server) sendVerificationCodeForEmailChange(email simple.Email, code simple.VerificationCode) (re *jrm1.RpcError) {
	var subject = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplateForRegVCode.ToString(), srv.settings.SystemSettings.SiteName))
	var msg = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForEmailChange.ToString(), code))
	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}
	return srv.sendEmailMessage(params)
}

func (srv *Server) sendVerificationCodeForPwdChange(email simple.Email, code simple.VerificationCode) (re *jrm1.RpcError) {
	var subject = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplateForRegVCode.ToString(), srv.settings.SystemSettings.SiteName))
	var msg = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForPwdChange.ToString(), code))
	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}
	return srv.sendEmailMessage(params)
}

func (srv *Server) sendGreetingAfterReg(email simple.Email) (re *jrm1.RpcError) {
	var subject = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplateForReg.ToString(), srv.settings.SystemSettings.SiteName))
	var msg = base2.Text(fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForReg.ToString(), srv.settings.SystemSettings.SiteName))
	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}
	return srv.sendEmailMessage(params)
}
