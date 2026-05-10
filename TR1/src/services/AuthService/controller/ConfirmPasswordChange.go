package c

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/vault-thirteen/BytePackedPassword"
	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
)

func (c *Controller) ConfirmPasswordChange(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ConfirmPasswordChangeParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ConfirmPasswordChangeResult
	r, re = c.confirmPasswordChange(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) confirmPasswordChange(p *rm.ConfirmPasswordChangeParams) (result *rm.ConfirmPasswordChangeResult, re *jrm1.RpcError) {
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
	var pcr = &cm.PasswordChangeRequest{
		RequestId: p.RequestId,
	}
	err = dbC.FindPasswordChangeRequest(pcr)
	if err != nil {
		return nil, c.databaseError(err)
	}

	// Check for fraud.
	{
		if !bytes.Equal(pcr.UserIPAB, p.Auth.UserIPAB) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
		if pcr.UserId != userWithSession.Id {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
	}

	// Check captcha & verification codes.
	{
		var isCorrect bool
		isCorrect, re = c.checkCaptcha(pcr.CaptchaId, p.CaptchaAnswer)
		if re != nil {
			return nil, re
		}
		if !isCorrect {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsWrong, rme.Msg_CaptchaAnswerIsWrong, nil)
		}

		if pcr.VerificationCode != p.VerificationCode {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
	}

	var userWithPassword *cm.User

	// Check password.
	{
		userWithPassword = &cm.User{Id: pcr.UserId}
		err = dbC.GetUserByIdAbleToLogIn(userWithPassword)
		if err != nil {
			return nil, c.databaseError(err)
		}

		var pwdRunes []rune
		pwdRunes, err = bpp.UnpackBytes(userWithPassword.Password.Bytes)
		if err != nil {
			c.logError(err)
			return nil, jrm1.NewRpcErrorByUser(rme.Code_BPP, fmt.Sprintf(rme.MsgF_BPP, err.Error()), nil)
		}

		var ok bool
		ok, err = bpp.CheckHashKey(string(pwdRunes), pcr.AuthData, p.AuthData)
		if err != nil {
			c.logError(err)
			return nil, jrm1.NewRpcErrorByUser(rme.Code_BPP, fmt.Sprintf(rme.MsgF_BPP, err.Error()), nil)
		}

		if !ok {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordIsWrong, rme.Msg_PasswordIsWrong, nil)
		}
	}

	// Confirm password change.
	{
		err = dbC.SaveUserPassword(userWithPassword, pcr.NewPassword)
		if err != nil {
			return nil, c.databaseError(err)
		}

		err = dbC.DeletePasswordChangeRequest(pcr)
		if err != nil {
			return nil, c.databaseError(err)
		}

		re = c.logOutUserByAction(userWithSession.Id, userWithSession.Session.Id)
		if re != nil {
			return nil, re
		}

		// Result.
		result = &rm.ConfirmPasswordChangeResult{
			Success:     rm.Success{OK: true},
			IsLoggedOut: true,
		}

		return result, nil
	}
}
