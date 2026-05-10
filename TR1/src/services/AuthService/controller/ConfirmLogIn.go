package c

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vault-thirteen/BytePackedPassword"
	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

func (c *Controller) ConfirmLogIn(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ConfirmLogInParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ConfirmLogInResult
	r, re = c.confirmLogIn(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) confirmLogIn(p *rm.ConfirmLogInParams) (result *rm.ConfirmLogInResult, re *jrm1.RpcError) {
	// Access check.
	{
		re = c.mustBeNoAuthToken(p.Auth)
		if re != nil {
			return nil, re
		}
	}

	// Check parameters.
	{
		if len(p.RequestId) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_RequestIdIsNotSet, rme.Msg_RequestIdIsNotSet, nil)
		}
		if len(p.CaptchaAnswer) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsNotSet, rme.Msg_CaptchaAnswerIsNotSet, nil)
		}
		if len(p.VerificationCode) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_VerificationCodeIsNotSet, rme.Msg_VerificationCodeIsNotSet, nil)
		}
		if len(p.AuthData) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_AuthDataIsNotSet, rme.Msg_AuthDataIsNotSet, nil)
		}
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())
	var lir = &cm.LogInRequest{
		RequestId: p.RequestId,
	}
	err = dbC.FindLogInRequest(lir)
	if err != nil {
		return nil, c.databaseError(err)
	}

	// Check for fraud.
	{
		if !bytes.Equal(lir.UserIPAB, p.Auth.UserIPAB) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
	}

	// Check captcha & verification code.
	{
		var isCorrect bool
		isCorrect, re = c.checkCaptcha(lir.CaptchaId, p.CaptchaAnswer)
		if re != nil {
			return nil, re
		}
		if !isCorrect {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsWrong, rme.Msg_CaptchaAnswerIsWrong, nil)
		}

		if lir.VerificationCode != p.VerificationCode {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_VerificationCodeIsWrong, rme.Msg_VerificationCodeIsWrong, nil)
		}
	}

	var user *cm.User

	// Check password.
	{
		user = &cm.User{Id: lir.UserId}
		err = dbC.GetUserByIdAbleToLogIn(user)
		if err != nil {
			return nil, c.databaseError(err)
		}

		var pwdRunes []rune
		pwdRunes, err = bpp.UnpackBytes(user.Password.Bytes)
		if err != nil {
			c.logError(err)
			return nil, jrm1.NewRpcErrorByUser(rme.Code_BPP, fmt.Sprintf(rme.MsgF_BPP, err.Error()), nil)
		}

		var ok bool
		ok, err = bpp.CheckHashKey(string(pwdRunes), lir.AuthData, p.AuthData)
		if err != nil {
			c.logError(err)
			return nil, jrm1.NewRpcErrorByUser(rme.Code_BPP, fmt.Sprintf(rme.MsgF_BPP, err.Error()), nil)
		}

		if !ok {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordIsWrong, rme.Msg_PasswordIsWrong, nil)
		}
	}

	// Confirm user logging in.
	{
		// Create session.
		var session = &cm.Session{
			UserId:   user.Id,
			UserIPAB: p.Auth.UserIPAB,
		}

		err = dbC.CreateSession(session)
		if err != nil {
			return nil, c.databaseError(err)
		}

		sessionMaxDurationSec := c.far.systemSettings.GetParameterAsInt(ccp.SessionMaxDuration)
		expirationTime := time.Now().Add(time.Duration(sessionMaxDurationSec) * time.Second)

		var token string
		token, err = c.far.jwtkm.MakeJWToken(user.Id, session.Id, expirationTime)
		if err != nil {
			c.logError(err)
			return nil, jrm1.NewRpcErrorByUser(rme.Code_JWT, fmt.Sprintf(rme.MsgF_JWT, err.Error()), nil)
		}

		err = dbC.DeleteLogInRequest(lir)
		if err != nil {
			return nil, c.databaseError(err)
		}

		// Journaling.
		logEvent := cm.NewLogEvent(cm.LogEvent_Type_LogIn, user.Id, p.Auth.UserIPAB, nil)

		err = dbC.CreateLogEvent(logEvent)
		if err != nil {
			return nil, c.databaseError(err)
		}

		// Result.
		result = &rm.ConfirmLogInResult{
			Success:    rm.Success{OK: true},
			IsTokenSet: true,
			Token:      token,
		}

		return result, nil
	}
}
