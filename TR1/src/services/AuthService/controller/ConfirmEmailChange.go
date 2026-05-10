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

func (c *Controller) ConfirmEmailChange(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ConfirmEmailChangeParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ConfirmEmailChangeResult
	r, re = c.confirmEmailChange(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) confirmEmailChange(p *rm.ConfirmEmailChangeParams) (result *rm.ConfirmEmailChangeResult, re *jrm1.RpcError) {
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
		if len(p.VerificationCodeA) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_VerificationCodeIsNotSet, rme.Msg_VerificationCodeIsNotSet, nil)
		}
		if len(p.VerificationCodeB) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_VerificationCodeIsNotSet, rme.Msg_VerificationCodeIsNotSet, nil)
		}
		if len(p.AuthData) == 0 {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_AuthDataIsNotSet, rme.Msg_AuthDataIsNotSet, nil)
		}
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())
	var ecr = &cm.EmailChangeRequest{
		RequestId: p.RequestId,
	}
	err = dbC.FindEmailChangeRequest(ecr)
	if err != nil {
		return nil, c.databaseError(err)
	}

	// Check for fraud.
	{
		if !bytes.Equal(ecr.UserIPAB, p.Auth.UserIPAB) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
		if ecr.UserId != userWithSession.Id {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
	}

	// Check captcha & verification codes.
	{
		var isCorrect bool
		isCorrect, re = c.checkCaptcha(ecr.CaptchaId, p.CaptchaAnswer)
		if re != nil {
			return nil, re
		}
		if !isCorrect {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsWrong, rme.Msg_CaptchaAnswerIsWrong, nil)
		}

		if ecr.VerificationCodeA != p.VerificationCodeA {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
		if ecr.VerificationCodeB != p.VerificationCodeB {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
	}

	// Check password.
	{
		userWithPassword := &cm.User{Id: ecr.UserId}
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
		ok, err = bpp.CheckHashKey(string(pwdRunes), ecr.AuthData, p.AuthData)
		if err != nil {
			c.logError(err)
			return nil, jrm1.NewRpcErrorByUser(rme.Code_BPP, fmt.Sprintf(rme.MsgF_BPP, err.Error()), nil)
		}

		if !ok {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_PasswordIsWrong, rme.Msg_PasswordIsWrong, nil)
		}
	}

	// Confirm e-mail change.
	{
		err = dbC.SaveUserEmail(userWithSession, ecr.NewEmail)
		if err != nil {
			return nil, c.databaseError(err)
		}

		err = dbC.DeleteEmailChangeRequest(ecr)
		if err != nil {
			return nil, c.databaseError(err)
		}

		re = c.logOutUserByAction(userWithSession.Id, userWithSession.Session.Id)
		if re != nil {
			return nil, re
		}

		// Result.
		result = &rm.ConfirmEmailChangeResult{
			Success:     rm.Success{OK: true},
			IsLoggedOut: true,
		}

		return result, nil
	}
}
