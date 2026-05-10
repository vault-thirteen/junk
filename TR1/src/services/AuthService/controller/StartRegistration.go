package c

import (
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

func (c *Controller) StartRegistration(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.StartRegistrationParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.StartRegistrationResult
	r, re = c.startRegistration(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) startRegistration(p *rm.StartRegistrationParams) (result *rm.StartRegistrationResult, re *jrm1.RpcError) {
	// Access check.
	{
		re = c.mustBeNoAuthToken(p.Auth)
		if re != nil {
			return nil, re
		}
	}

	// Check parameters.
	{
		if p.User == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserIsNotSet, rme.Msg_UserIsNotSet, nil)
		}
		if len(p.User.Name) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_NameIsNotSet, rme.Msg_NameIsNotSet, nil)
		}
		if len(p.User.Email) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_EmailIsNotSet, rme.Msg_EmailIsNotSet, nil)
		}
		if p.User.Password == nil {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordIsNotSet, rme.Msg_PasswordIsNotSet, nil)
		}
		if len(p.User.Password.Text) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordIsNotSet, rme.Msg_PasswordIsNotSet, nil)
		}
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())

	// Check for existing user with name & e-mail.
	{
		var isFree bool
		isFree, err = dbC.IsUserNameFree(p.User.Name)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if !isFree {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserNameIsUsed, rme.Msg_UserNameIsUsed, p.User.Name)
		}

		isFree, err = dbC.IsUserEmailFree(p.User.Email)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if !isFree {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserEmailIsUsed, rme.Msg_UserEmailIsUsed, p.User.Email)
		}
	}

	// Check for existing registration request with name & e-mail.
	{
		var exists bool
		exists, err = dbC.ExistsRegistrationRequestWithUserName(p.User.Name)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if exists {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_RegistrationRequestWithUserNameExists, rme.Msg_RegistrationRequestWithUserNameExists, p.User.Name)
		}

		exists, err = dbC.ExistsRegistrationRequestWithUserEmail(p.User.Email)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if exists {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_RegistrationRequestWithUserEmailExists, rme.Msg_RegistrationRequestWithUserEmailExists, p.User.Email)
		}
	}

	// Other checks.
	{
		if !cm.IsUserEmailValid(p.User.Email) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserEmailIsInvalid, rme.Msg_UserEmailIsInvalid, p.User.Email)
		}

		if len([]byte(p.User.Name)) > c.far.systemSettings.GetParameterAsInt(ccp.UserNameMaxLenInBytes) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserNameIsTooLong, rme.Msg_UserNameIsTooLong, p.User.Name)
		}

		if len([]byte(p.User.Password.Text)) > c.far.systemSettings.GetParameterAsInt(ccp.UserPasswordMaxLenInBytes) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserPasswordIsTooLong, rme.Msg_UserPasswordIsTooLong, nil)
		}

		ok := cm.IsUserPasswordAllowed(p.User.Password.Text)
		if !ok {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserPasswordIsNotAllowed, rme.Msg_UserPasswordIsNotAllowed, nil)
		}
	}

	// Start user registration.
	{
		var requestId *string
		requestId, re = c.createRequestId()
		if re != nil {
			return nil, re
		}

		var captchaData *rm.CreateCaptchaResult
		captchaData, re = c.createCaptcha()
		if re != nil {
			return nil, re
		}

		var verificationCode *string
		verificationCode, re = c.createVerificationCode()
		if re != nil {
			return nil, re
		}

		re = c.sendVerificationCode_Reg(p.User.Email, *verificationCode)
		if re != nil {
			return nil, re
		}

		var rr = cm.RegistrationRequest{
			UserName:     p.User.Name,
			UserEmail:    p.User.Email,
			UserPassword: p.User.Password.Text,

			RequestId:        *requestId,
			UserIPAB:         p.Auth.UserIPAB,
			CaptchaId:        captchaData.TaskId,
			VerificationCode: *verificationCode,
		}
		err = dbC.CreateRegistrationRequest(rr)
		if err != nil {
			return nil, c.databaseError(err)
		}

		result = &rm.StartRegistrationResult{
			RequestId: *requestId,
			CaptchaId: captchaData.TaskId,
		}

		return result, nil
	}
}
