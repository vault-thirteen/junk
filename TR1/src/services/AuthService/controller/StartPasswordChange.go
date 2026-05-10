package c

import (
	"encoding/json"
	"fmt"

	"github.com/vault-thirteen/BytePackedPassword"
	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	rme "github.com/vault-thirteen/TR1/src/models/rpc/error"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

func (c *Controller) StartPasswordChange(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.StartPasswordChangeParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.StartPasswordChangeResult
	r, re = c.startPasswordChange(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) startPasswordChange(p *rm.StartPasswordChangeParams) (result *rm.StartPasswordChangeResult, re *jrm1.RpcError) {
	var userWithSession *cm.User

	// Access check.
	{
		userWithSession, re = c.mustBeAnAuthToken(p.Auth)
		if re != nil {
			return nil, re
		}
	}

	// Check parameters.
	{
		if len(p.NewPassword) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordIsNotSet, rme.Msg_PasswordIsNotSet, nil)
		}
		if len(p.NewPassword2) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordIsNotSet, rme.Msg_PasswordIsNotSet, nil)
		}
		if p.NewPassword != p.NewPassword2 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordIsNotSet, rme.Msg_PasswordIsNotSet, nil)
		}
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())

	// Check for existing request.
	{
		var exists bool
		exists, err = dbC.ExistsPasswordChangeRequestWithUserId(userWithSession)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if exists {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordChangeRequestWithUserIdExists, rme.Msg_PasswordChangeRequestWithUserIdExists, userWithSession.Id)
		}
	}

	// Other checks.
	{
		if len([]byte(p.NewPassword)) > c.far.systemSettings.GetParameterAsInt(ccp.UserPasswordMaxLenInBytes) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserPasswordIsTooLong, rme.Msg_UserPasswordIsTooLong, nil)
		}

		ok := cm.IsUserPasswordAllowed(p.NewPassword)
		if !ok {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserPasswordIsNotAllowed, rme.Msg_UserPasswordIsNotAllowed, nil)
		}
	}

	// Start changing user password.
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

		// Verification code.
		var verificationCode *string
		verificationCode, re = c.createVerificationCode()
		if re != nil {
			return nil, re
		}

		re = c.sendVerificationCode_PwdChange(userWithSession.Email, *verificationCode)
		if re != nil {
			return nil, re
		}

		var pwdSalt []byte
		pwdSalt, err = bpp.GenerateRandomSalt()
		if err != nil {
			c.logError(err)
			return nil, jrm1.NewRpcErrorByUser(rme.Code_BPP, fmt.Sprintf(rme.MsgF_BPP, err.Error()), nil)
		}

		var pcr = cm.PasswordChangeRequest{
			NewPassword:      p.NewPassword,
			UserId:           userWithSession.Id,
			RequestId:        *requestId,
			UserIPAB:         p.Auth.UserIPAB,
			CaptchaId:        captchaData.TaskId,
			VerificationCode: *verificationCode,
			AuthData:         pwdSalt,
		}
		err = dbC.CreatePasswordChangeRequest(pcr)
		if err != nil {
			return nil, c.databaseError(err)
		}

		result = &rm.StartPasswordChangeResult{
			RequestId: *requestId,
			CaptchaId: captchaData.TaskId,
			AuthData:  pwdSalt,
		}

		return result, nil
	}
}
