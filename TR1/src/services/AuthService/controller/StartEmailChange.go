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
)

func (c *Controller) StartEmailChange(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.StartEmailChangeParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.StartEmailChangeResult
	r, re = c.startEmailChange(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) startEmailChange(p *rm.StartEmailChangeParams) (result *rm.StartEmailChangeResult, re *jrm1.RpcError) {
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
		if len(p.NewEmail) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_EmailIsNotSet, rme.Msg_EmailIsNotSet, nil)
		}
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())

	// Check for existing user with e-mail.
	{
		var isFree bool
		isFree, err = dbC.IsUserEmailFree(p.NewEmail)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if !isFree {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserEmailIsUsed, rme.Msg_UserEmailIsUsed, p.NewEmail)
		}
	}

	// Check for existing request with e-mail.
	{
		var exists bool
		exists, err = dbC.ExistsEmailChangeRequestWithNewEmail(p.NewEmail)
		if err != nil {
			return nil, c.databaseError(err)
		}
		if exists {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_EmailChangeRequestWithNewEmailExists, rme.Msg_EmailChangeRequestWithNewEmailExists, p.NewEmail)
		}
	}

	// Other checks.
	{
		if !cm.IsUserEmailValid(p.NewEmail) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_UserEmailIsInvalid, rme.Msg_UserEmailIsInvalid, p.NewEmail)
		}
	}

	// Start changing e-mail.
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

		// Verification code for old e-mail address.
		var verificationCodeA *string
		verificationCodeA, re = c.createVerificationCode()
		if re != nil {
			return nil, re
		}

		re = c.sendVerificationCode_EmailChange(userWithSession.Email, *verificationCodeA)
		if re != nil {
			return nil, re
		}

		// Verification code for new e-mail address.
		var verificationCodeB *string
		verificationCodeB, re = c.createVerificationCode()
		if re != nil {
			return nil, re
		}

		re = c.sendVerificationCode_EmailChange(p.NewEmail, *verificationCodeB)
		if re != nil {
			return nil, re
		}

		var pwdSalt []byte
		pwdSalt, err = bpp.GenerateRandomSalt()
		if err != nil {
			c.logError(err)
			return nil, jrm1.NewRpcErrorByUser(rme.Code_BPP, fmt.Sprintf(rme.MsgF_BPP, err.Error()), nil)
		}

		var ecr = cm.EmailChangeRequest{
			NewEmail:          p.NewEmail,
			UserId:            userWithSession.Id,
			RequestId:         *requestId,
			UserIPAB:          p.Auth.UserIPAB,
			CaptchaId:         captchaData.TaskId,
			VerificationCodeA: *verificationCodeA,
			VerificationCodeB: *verificationCodeB,
			AuthData:          pwdSalt,
		}
		err = dbC.CreateEmailChangeRequest(ecr)
		if err != nil {
			return nil, c.databaseError(err)
		}

		result = &rm.StartEmailChangeResult{
			RequestId: *requestId,
			CaptchaId: captchaData.TaskId,
			AuthData:  pwdSalt,
		}

		return result, nil
	}
}
