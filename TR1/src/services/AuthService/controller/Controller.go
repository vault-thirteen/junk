package c

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/components/CaptchaComponent"
	"github.com/vault-thirteen/TR1/src/components/DatabaseComponent"
	"github.com/vault-thirteen/TR1/src/components/JwtManagerComponent"
	"github.com/vault-thirteen/TR1/src/components/RequestIdGeneratorComponent"
	"github.com/vault-thirteen/TR1/src/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/components/VerificationCodeGeneratorComponent"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/libraries/net"
	"github.com/vault-thirteen/TR1/src/libraries/scheduler"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	"github.com/vault-thirteen/auxie/number"
	"gorm.io/gorm"
)

// List of component indices of the controller must be synchronised with the
// order of components used in the application's constructor.
const (
	ComponentIndex_ConsoleComponent                   = 0
	ComponentIndex_ErrorListenerComponent             = 1
	ComponentIndex_DatabaseComponent                  = 2
	ComponentIndex_JwtManagerComponent                = 3
	ComponentIndex_RequestIdGeneratorComponent        = 4
	ComponentIndex_RpcClientComponent                 = 5
	ComponentIndex_VerificationCodeGeneratorComponent = 6
	ComponentIndex_RpcServerComponent                 = 7
	ComponentIndex_SchedulerComponent                 = 8
)

type Controller struct {
	cfg        *cm.Configuration
	errorsChan *chan error
	service    *cm.Service
	far        ControllerFastAccessRegistry
}

func NewController() (c *Controller) {
	errorsChan := make(chan error, 1)

	return &Controller{
		errorsChan: &errorsChan,
	}
}

func (c *Controller) GetRpcFunctions() []jrm1.RpcFunction {
	return []jrm1.RpcFunction{
		// Ping.
		c.Ping,

		// User registration.
		c.StartRegistration,
		c.ConfirmRegistration,
		c.ListRegistrationRequestsRFA,
		c.ApproveRegistrationRequestRFA,
		c.RejectRegistrationRequestRFA,

		// Logging in and out.
		c.StartLogIn,
		c.ConfirmLogIn,
		c.StartLogOut,
		c.ConfirmLogOut,
		c.IsUserLoggedIn,
		c.LogUserOutA,

		// User properties.
		c.GetUserName,
		c.GetUserRoles,
		c.GetSelfRoles,
		c.GetUserParameters,
		c.SetUserRoleAuthor,
		c.SetUserRoleWriter,
		c.SetUserRoleReader,
		c.BanUser,
		c.UnbanUser,
		c.StartEmailChange,
		c.ConfirmEmailChange,
		c.StartPasswordChange,
		c.ConfirmPasswordChange,

		// Various actions.
		c.ListUsers,
		c.ListUserSessions,
		c.GetUserSession,
	}
}

func (c *Controller) GetScheduledFunctions() []sch.ScheduledFn {
	return []sch.ScheduledFn{
		c.RemoveOutdatedRegistrationRequests,
		c.RemoveOutdatedLogInRequests,
		c.RemoveOutdatedLogOutRequests,
		c.RemoveOutdatedEmailChangeRequests,
		c.RemoveOutdatedPasswordChangeRequests,
		c.RemoveOutdatedSessions,
	}
}

func (c *Controller) GetErrorsChan() (errorsChan *chan error) {
	return c.errorsChan
}

func (c *Controller) LinkWithService(service interfaces.IService) (err error) {
	c.cfg = (service.GetConfiguration()).(*cm.Configuration)
	c.service = service.(*cm.Service)
	c.initFAR()

	err = c.prepareDb()
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) initFAR() {
	c.far = ControllerFastAccessRegistry{}

	c.far.systemSettings = c.cfg.GetComponent(cm.Component_System, cm.Protocol_None)
	c.far.messageSettings = c.cfg.GetComponent(cm.Component_Message, cm.Protocol_None)
	c.far.roleSettings = c.cfg.GetComponent(cm.Component_Role, cm.Protocol_None)

	c.far.rcc = rcc.FromAny(c.service.GetComponentByIndex(ComponentIndex_RpcClientComponent))

	c.far.rcsServiceClient = c.far.rcc.GetClientMap()[rm.ServiceShortName_RCS]
	c.far.mailerServiceClient = c.far.rcc.GetClientMap()[rm.ServiceShortName_Mailer]

	c.far.pageSize = c.far.systemSettings.GetParameterAsInt(ccp.PageSize)

	c.far.dbc = dc.FromAny(c.service.GetComponentByIndex(ComponentIndex_DatabaseComponent))
	c.far.db = c.far.dbc.GetGormDb()
	c.far.ridgc = rigc.FromAny(c.service.GetComponentByIndex(ComponentIndex_RequestIdGeneratorComponent))
	c.far.ridg = c.far.ridgc.GetRidg()
	c.far.vcgc = vcgc.FromAny(c.service.GetComponentByIndex(ComponentIndex_VerificationCodeGeneratorComponent))
	c.far.vcg = c.far.vcgc.GetVcg()
	c.far.jmc = jmc.FromAny(c.service.GetComponentByIndex(ComponentIndex_JwtManagerComponent))
	c.far.jwtkm = c.far.jmc.GetKeyMaker()
}

func (c *Controller) prepareDb() (err error) {
	db := c.GetDb()

	if c.far.systemSettings.GetParameterAsBool(ccp.IsDatabaseInitialisationUsed) {
		classesToInit := []any{
			&cm.EmailChangeRequest{},
			&cm.LogEvent{},
			&cm.LogInRequest{},
			&cm.LogOutRequest{},
			&cm.Password{},
			&cm.PasswordChangeRequest{},
			&cm.RegistrationRequest{},
			&cm.Session{},
			&cm.User{},
		}

		for _, cti := range classesToInit {
			err = db.AutoMigrate(cti)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Controller) GetDb() (gormDb *gorm.DB) {
	return c.far.db
}

func (c *Controller) logError(err error) {
	if err == nil {
		return
	}

	if c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode) {
		log.Println(err)
	}
}
func (c *Controller) databaseError(err error) (re *jrm1.RpcError) {
	c.processDatabaseError(err)
	return jrm1.NewRpcErrorByUser(rme.Code_Database, rme.Msg_Database, err)
}
func (c *Controller) processDatabaseError(err error) {
	if err == nil {
		return
	}

	if rme.IsNetworkError(err) {
		log.Println(fmt.Sprintf(rme.ErrF_DatabaseNetwork, err.Error()))
		*(c.errorsChan) <- err
	} else {
		c.logError(err)
	}

	return
}

func (c *Controller) mustBeNoAuthToken(auth *rm.Auth) (re *jrm1.RpcError) {
	re = c.mustBeAuthUserIPA(auth)
	if re != nil {
		return re
	}

	if len(auth.Token) > 0 {
		return jrm1.NewRpcErrorByUser(rme.Code_Permission, rme.Msg_Permission, nil)
	}

	return nil
}
func (c *Controller) mustBeAuthUserIPA(auth *rm.Auth) (re *jrm1.RpcError) {
	if (auth == nil) || (len(auth.UserIPA) == 0) {
		return jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
	}

	var err error
	auth.UserIPAB, err = net.ParseIPA(auth.UserIPA)
	if err != nil {
		c.logError(err)
		return jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
	}

	return nil
}
func (c *Controller) mustBeAnAuthToken(auth *rm.Auth) (userWithSession *cm.User, re *jrm1.RpcError) {
	re = c.mustBeAuthUserIPA(auth)
	if re != nil {
		return nil, re
	}

	if len(auth.Token) == 0 {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_NotAuthorised, rme.Msg_NotAuthorised, nil)
	}

	userWithSession, re = c.getUserWithSessionByAuthToken(auth.Token)
	if re != nil {
		return nil, re
	}

	if bytes.Compare(auth.UserIPAB, userWithSession.Session.UserIPAB) != 0 {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
	}

	return userWithSession, nil
}
func (c *Controller) getUserWithSessionByAuthToken(authToken string) (userWithSession *cm.User, re *jrm1.RpcError) {
	dbC := dbc.NewDbController(c.GetDb())

	var userId, sessionId int
	var err error
	userId, sessionId, err = c.far.jwtkm.ValidateToken(authToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			re = c.logOutUserByTimeout(userId, sessionId)
			if re != nil {
				return nil, re
			}

			return nil, jrm1.NewRpcErrorByUser(rme.Code_TokenIsExpired, rme.Msg_TokenIsExpired, nil)
		}

		c.logError(err)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
	}

	userWithSession = &cm.User{Id: userId}
	err = dbC.GetUserWithSessionByIdAbleToLogIn(userWithSession)
	if err != nil {
		return nil, c.databaseError(err)
	}

	if (userWithSession.Session == nil) || (userWithSession.Session.Id != sessionId) {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_SessionIsNotFound, rme.Msg_SessionIsNotFound, nil)
	}

	c.attachUserSpecialRoles(userWithSession)

	return userWithSession, nil
}
func (c *Controller) attachUserSpecialRoles(user *cm.User) {
	user.Roles.IsModerator = c.isUserModerator(user.Id)
	user.Roles.IsAdministrator = c.isUserAdministrator(user.Id)
}
func (c *Controller) attachUsersSpecialRoles(users []cm.User) {
	administrators := c.far.roleSettings.GetParameterAsInts(ccp.Administrator)
	moderators := c.far.roleSettings.GetParameterAsInts(ccp.Moderator)

	for _, user := range users {
		for _, id := range administrators {
			if id == user.Id {
				user.Roles.IsAdministrator = true
			}
		}

		for _, id := range moderators {
			if id == user.Id {
				user.Roles.IsModerator = true
			}
		}
	}
}
func (c *Controller) attachUserSessionsSpecialRoles(userSessions []cm.Session) {
	administrators := c.far.roleSettings.GetParameterAsInts(ccp.Administrator)
	moderators := c.far.roleSettings.GetParameterAsInts(ccp.Moderator)

	for _, us := range userSessions {
		for _, id := range administrators {
			if id == us.User.Id {
				us.User.Roles.IsAdministrator = true
			}
		}

		for _, id := range moderators {
			if id == us.User.Id {
				us.User.Roles.IsModerator = true
			}
		}
	}
}
func (c *Controller) isUserAdministrator(userId int) (isAdministrator bool) {
	for _, id := range c.far.roleSettings.GetParameterAsInts(ccp.Administrator) {
		if id == userId {
			return true
		}
	}
	return false
}
func (c *Controller) isUserModerator(userId int) (isModerator bool) {
	for _, id := range c.far.roleSettings.GetParameterAsInts(ccp.Moderator) {
		if id == userId {
			return true
		}
	}
	return false
}

func (c *Controller) registerUser(rr *cm.RegistrationRequest) (re *jrm1.RpcError) {
	dbC := dbc.NewDbController(c.GetDb())

	var user = &cm.User{
		Name:  rr.UserName,
		Email: rr.UserEmail,
		Roles: &cm.Roles{
			CanLogIn: true,
			CanRead:  true,
		},
		RegTime: time.Now(),
	}

	err := dbC.CreateUser(user, rr.UserPassword)
	if err != nil {
		return c.databaseError(err)
	}

	err = dbC.DeleteRegistrationRequestRFA(rr)
	if err != nil {
		return c.databaseError(err)
	}

	return nil
}
func (c *Controller) logOutUserBySelf(userId int, sessionId int) (re *jrm1.RpcError) {
	return c.logOutUser(userId, sessionId, cm.LogEvent_Type_LogOutBySelf, nil)
}
func (c *Controller) logOutUserByTimeout(userId int, sessionId int) (re *jrm1.RpcError) {
	return c.logOutUser(userId, sessionId, cm.LogEvent_Type_LogOutByTimeout, nil)
}
func (c *Controller) logOutUserByAction(userId int, sessionId int) (re *jrm1.RpcError) {
	return c.logOutUser(userId, sessionId, cm.LogEvent_Type_LogOutByAction, nil)
}
func (c *Controller) logOutUserByAdministrator(userId int, sessionId int, adminId *int) (re *jrm1.RpcError) {
	return c.logOutUser(userId, sessionId, cm.LogEvent_Type_LogOutByAdministrator, adminId)
}
func (c *Controller) logOutUser(userId int, sessionId int, logEventType int, adminId *int) (re *jrm1.RpcError) {
	dbC := dbc.NewDbController(c.GetDb())
	user := &cm.User{Id: userId}

	err := dbC.GetUserWithSessionByIdAbleToLogIn(user)
	if err != nil {
		return c.databaseError(err)
	}

	if (user.Session == nil) || (user.Session.Id != sessionId) {
		return jrm1.NewRpcErrorByUser(rme.Code_SessionIsNotFound, rme.Msg_SessionIsNotFound, nil)
	}

	// Delete session.
	err = dbC.DeleteSession(user.Session)
	if err != nil {
		return c.databaseError(err)
	}

	// Journaling.
	logEvent := cm.NewLogEvent(logEventType, user.Id, nil, adminId)

	err = dbC.CreateLogEvent(logEvent)
	if err != nil {
		return c.databaseError(err)
	}

	return nil
}

func (c *Controller) createRequestId() (rid *string, re *jrm1.RpcError) {
	var err error
	rid, err = c.far.ridg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_RequestIdGenerator, rme.Msg_RequestIdGenerator, nil)
	}

	return rid, nil
}
func (c *Controller) createCaptcha() (result *rm.CreateCaptchaResult, re *jrm1.RpcError) {
	var params = rm.CreateCaptchaParams{}
	result = new(rm.CreateCaptchaResult)

	var err error
	re, err = c.far.rcsServiceClient.MakeRequest(context.Background(), rm.Func_CreateCaptcha, params, result)
	if err != nil {
		c.logError(err)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_RPCCall, rme.Msg_RPCCall, nil)
	}
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaError, fmt.Sprintf(rme.MsgF_CaptchaError, re.AsError().Error()), nil)
	}

	if result.IsImageDataReturned {
		// We do not return an image in a JSON message.
		err = errors.New(cc.Err_UnexpectedResponse)
		return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaError, fmt.Sprintf(rme.MsgF_CaptchaError, err.Error()), nil)
	}

	return result, nil
}
func (c *Controller) createVerificationCode() (vc *string, re *jrm1.RpcError) {
	var err error
	var s *string
	s, err = c.far.vcg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(rme.Code_VerificationCodeGenerator, rme.Msg_VerificationCodeGenerator, nil)
	}

	return s, nil
}
func (c *Controller) sendVerificationCode_Reg(email string, vCode string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.SubjectTemplate_VC), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var msg = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.BodyTemplate_Reg), vCode)
	var params = rm.SendEmailMessageParams{Recipient: email, Subject: subject, Message: msg}
	return c.sendEmailMessage(params)
}
func (c *Controller) sendMessage_RegRFA(email string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.SubjectTemplate_Reg), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var msg = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.BodyTemplate_RegRFA), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var params = rm.SendEmailMessageParams{Recipient: email, Subject: subject, Message: msg}
	return c.sendEmailMessage(params)
}
func (c *Controller) sendMessage_RegApproved(email string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.SubjectTemplate_Reg), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var msg = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.BodyTemplate_RegApproved), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var params = rm.SendEmailMessageParams{Recipient: email, Subject: subject, Message: msg}
	return c.sendEmailMessage(params)
}
func (c *Controller) sendVerificationCode_LogIn(email string, vCode string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.SubjectTemplate_VC), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var msg = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.BodyTemplate_LogIn), vCode)
	var params = rm.SendEmailMessageParams{Recipient: email, Subject: subject, Message: msg}
	return c.sendEmailMessage(params)
}
func (c *Controller) sendVerificationCode_PwdChange(email string, vCode string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.SubjectTemplate_VC), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var msg = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.BodyTemplate_PwdChange), vCode)
	var params = rm.SendEmailMessageParams{Recipient: email, Subject: subject, Message: msg}
	return c.sendEmailMessage(params)
}
func (c *Controller) sendVerificationCode_EmailChange(email string, vCode string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.SubjectTemplate_VC), c.far.systemSettings.GetParameterAsString(ccp.SiteName))
	var msg = fmt.Sprintf(c.far.messageSettings.GetParameterAsString(ccp.BodyTemplate_EmailChange), vCode)
	var params = rm.SendEmailMessageParams{Recipient: email, Subject: subject, Message: msg}
	return c.sendEmailMessage(params)
}
func (c *Controller) sendEmailMessage(params rm.SendEmailMessageParams) (re *jrm1.RpcError) {
	var result = new(rm.SendEmailMessageResult)

	var err error
	re, err = c.far.mailerServiceClient.MakeRequest(context.Background(), rm.Func_SendEmailMessage, params, result)
	if err != nil {
		c.logError(err)
		return jrm1.NewRpcErrorByUser(rme.Code_RPCCall, rme.Msg_RPCCall, nil)
	}
	if re != nil {
		return jrm1.NewRpcErrorByUser(rme.Code_MailerError, fmt.Sprintf(rme.MsgF_MailerError, re.AsError().Error()), err)
	}

	return nil
}
func (c *Controller) checkCaptcha(captchaId string, answer string) (isCorrect bool, re *jrm1.RpcError) {
	n, err := number.ParseUint(answer)
	if err != nil {
		return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaError, fmt.Sprintf(rme.MsgF_CaptchaError, err.Error()), answer)
	}

	var params = rm.CheckCaptchaParams{
		TaskId: captchaId,
		Value:  n,
	}

	// Fool check.
	{
		if len(captchaId) == 0 {
			return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaTaskIdIsNotSet, rme.Msg_CaptchaTaskIdIsNotSet, nil)
		}
		if len(answer) == 0 {
			return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsNotSet, rme.Msg_CaptchaAnswerIsNotSet, nil)
		}

		params.Value, err = number.ParseUint(answer)
		if err != nil {
			return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsNotSet, rme.Msg_CaptchaAnswerIsNotSet, nil)
		}
	}

	var result = new(rm.CheckCaptchaResult)
	re, err = c.far.rcsServiceClient.MakeRequest(context.Background(), rm.Func_CheckCaptcha, params, result)
	if err != nil {
		c.logError(err)
		return false, jrm1.NewRpcErrorByUser(rme.Code_RPCCall, rme.Msg_RPCCall, nil)
	}
	if re != nil {
		return false, jrm1.NewRpcErrorByUser(rme.Code_CaptchaCheckError, fmt.Sprintf(rme.Msg_CaptchaCheckError, re.AsError().Error()), nil)
	}

	return result.IsSuccess, nil
}
