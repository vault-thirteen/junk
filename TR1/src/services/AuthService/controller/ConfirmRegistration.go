package c

import (
	"bytes"
	"encoding/json"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/dbc"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/error"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

func (c *Controller) ConfirmRegistration(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *rm.ConfirmRegistrationParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *rm.ConfirmRegistrationResult
	r, re = c.confirmRegistration(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (c *Controller) confirmRegistration(p *rm.ConfirmRegistrationParams) (result *rm.ConfirmRegistrationResult, re *jrm1.RpcError) {
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
	}

	var err error
	dbC := dbc.NewDbController(c.GetDb())
	var rr = &cm.RegistrationRequest{
		RequestId: p.RequestId,
	}
	err = dbC.FindRegistrationRequestNRFA(rr)
	if err != nil {
		return nil, c.databaseError(err)
	}

	// Check for fraud.
	{
		if !bytes.Equal(rr.UserIPAB, p.Auth.UserIPAB) {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
	}

	// Check captcha & verification code.
	{
		var isCorrect bool
		isCorrect, re = c.checkCaptcha(rr.CaptchaId, p.CaptchaAnswer)
		if re != nil {
			return nil, re
		}
		if !isCorrect {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_CaptchaAnswerIsWrong, rme.Msg_CaptchaAnswerIsWrong, nil)
		}

		if rr.VerificationCode != p.VerificationCode {
			return nil, jrm1.NewRpcErrorByUser(rme.Code_Authorisation, rme.Msg_Authorisation, nil)
		}
	}

	// Confirm user registration.
	{
		err = dbC.MarkRegistrationRequestAsReadyForApproval(rr)
		if err != nil {
			return nil, c.databaseError(err)
		}

		// If approval is required, do nothing.
		if c.far.systemSettings.GetParameterAsBool(ccp.IsAdminApprovalRequired) {
			re = c.sendMessage_RegRFA(rr.UserEmail)
			if re != nil {
				return nil, re
			}

			result = &rm.ConfirmRegistrationResult{
				Success:            rm.Success{OK: true},
				IsApprovalRequired: true,
			}

			return result, nil
		}

		// Register user.
		re = c.registerUser(rr)
		if re != nil {
			return nil, re
		}

		re = c.sendMessage_RegApproved(rr.UserEmail)
		if re != nil {
			return nil, re
		}

		result = &rm.ConfirmRegistrationResult{
			Success:            rm.Success{OK: true},
			IsApprovalRequired: false,
		}

		return result, nil
	}
}
