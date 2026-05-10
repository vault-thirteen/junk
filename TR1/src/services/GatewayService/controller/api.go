package c

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	rm "github.com/vault-thirteen/TR1/src/models/rpc"
	rme "github.com/vault-thirteen/TR1/src/models/rpc/error"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

func (c *Controller) initAPI() {
	c.httpStatusCodesByRpcErrorCode = rme.GetMapOfHttpStatusCodesByRpcErrorCodes()

	c.apiHandlers = map[string]rm.RequestHandler{

		// AuthService.
		rm.ApiFunctionName_ApproveRegistrationRequestRFA: c.ApproveRegistrationRequestRFA,
		rm.ApiFunctionName_BanUser:                       c.BanUser,
		rm.ApiFunctionName_ConfirmEmailChange:            c.ConfirmEmailChange,
		rm.ApiFunctionName_ConfirmLogIn:                  c.ConfirmLogIn,
		rm.ApiFunctionName_ConfirmLogOut:                 c.ConfirmLogOut,
		rm.ApiFunctionName_ConfirmPasswordChange:         c.ConfirmPasswordChange,
		rm.ApiFunctionName_ConfirmRegistration:           c.ConfirmRegistration,
		rm.ApiFunctionName_GetSelfRoles:                  c.GetSelfRoles,
		rm.ApiFunctionName_GetUserName:                   c.GetUserName,
		rm.ApiFunctionName_GetUserParameters:             c.GetUserParameters,
		rm.ApiFunctionName_GetUserRoles:                  c.GetUserRoles,
		rm.ApiFunctionName_GetUserSession:                c.GetUserSession,
		rm.ApiFunctionName_IsUserLoggedIn:                c.IsUserLoggedIn,
		rm.ApiFunctionName_ListRegistrationRequestsRFA:   c.ListRegistrationRequestsRFA,
		rm.ApiFunctionName_ListUsers:                     c.ListUsers,
		rm.ApiFunctionName_ListUserSessions:              c.ListUserSessions,
		rm.ApiFunctionName_LogUserOutA:                   c.LogUserOutA,
		rm.ApiFunctionName_RejectRegistrationRequestRFA:  c.RejectRegistrationRequestRFA,
		rm.ApiFunctionName_SetUserRoleAuthor:             c.SetUserRoleAuthor,
		rm.ApiFunctionName_SetUserRoleReader:             c.SetUserRoleReader,
		rm.ApiFunctionName_SetUserRoleWriter:             c.SetUserRoleWriter,
		rm.ApiFunctionName_StartEmailChange:              c.StartEmailChange,
		rm.ApiFunctionName_StartLogIn:                    c.StartLogIn,
		rm.ApiFunctionName_StartLogOut:                   c.StartLogOut,
		rm.ApiFunctionName_StartPasswordChange:           c.StartPasswordChange,
		rm.ApiFunctionName_StartRegistration:             c.StartRegistration,
		rm.ApiFunctionName_UnbanUser:                     c.UnbanUser,

		// MessageService.
		rm.ApiFunctionName_AddForum:            c.AddForum,
		rm.ApiFunctionName_AddMessage:          c.AddMessage,
		rm.ApiFunctionName_AddThread:           c.AddThread,
		rm.ApiFunctionName_ChangeForumName:     c.ChangeForumName,
		rm.ApiFunctionName_ChangeMessageText:   c.ChangeMessageText,
		rm.ApiFunctionName_ChangeMessageThread: c.ChangeMessageThread,
		rm.ApiFunctionName_ChangeThreadForum:   c.ChangeThreadForum,
		rm.ApiFunctionName_ChangeThreadName:    c.ChangeThreadName,
		rm.ApiFunctionName_DeleteForum:         c.DeleteForum,
		rm.ApiFunctionName_DeleteMessage:       c.DeleteMessage,
		rm.ApiFunctionName_DeleteThread:        c.DeleteThread,
		rm.ApiFunctionName_GetForum:            c.GetForum,
		rm.ApiFunctionName_GetMessage:          c.GetMessage,
		rm.ApiFunctionName_GetThread:           c.GetThread,
		rm.ApiFunctionName_ListForums:          c.ListForums,
		rm.ApiFunctionName_ListMessages:        c.ListMessages,
		rm.ApiFunctionName_ListThreads:         c.ListThreads,
		rm.ApiFunctionName_MoveForumDown:       c.MoveForumDown,
		rm.ApiFunctionName_MoveForumUp:         c.MoveForumUp,
	}
}

func (c *Controller) handleApiRequest(rw http.ResponseWriter, req *http.Request, clientIPA string) {
	if req.Method != http.MethodPost {
		c.respondMethodNotAllowed(rw)
		return
	}

	// Check accepted MIME types.
	ok, err := hh.CheckBrowserSupportForJson(req)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}
	if !ok {
		c.respondNotAcceptable(rw)
		return
	}

	var reqBody []byte
	reqBody, err = io.ReadAll(req.Body)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}

	// Get the requested function.
	var arwoa rm.RequestWithOnlyAction
	err = json.Unmarshal(reqBody, &arwoa)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	if (arwoa.Action == nil) ||
		(arwoa.Parameters == nil) {
		c.respondBadRequest(rw)
		return
	}

	var handler rm.RequestHandler
	handler, ok = c.apiHandlers[*arwoa.Action]
	if !ok {
		c.respondNotFound(rw)
		return
	}

	var token *string
	token, err = rm.GetToken(req)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	var ar = &rm.Request{
		Action:     arwoa.Action,
		Parameters: arwoa.Parameters,
		Authorisation: &rm.Auth{
			UserIPA: clientIPA,
		},
	}

	if token != nil {
		ar.Authorisation.Token = *token
	}

	handler(ar, req, rw)
	return
}

// AuthService.

func (c *Controller) ApproveRegistrationRequestRFA(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ApproveRegistrationRequestRFAParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ApproveRegistrationRequestRFAResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_ApproveRegistrationRequestRFA, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) BanUser(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.BanUserParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.BanUserResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_BanUser, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ConfirmEmailChange(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ConfirmEmailChangeParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ConfirmEmailChangeResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_ConfirmEmailChange, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	// Clear HTTP cookies.
	c.clearTokenCookie(rw)

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ConfirmLogIn(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ConfirmLogInParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ConfirmLogInResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_ConfirmLogIn, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	// Set HTTP cookies.
	if result.IsTokenSet {
		c.setTokenCookie(rw, result.Token)
		result.Token = ""
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ConfirmLogOut(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ConfirmLogOutParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ConfirmLogOutResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_ConfirmLogOut, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	// Clear HTTP cookies.
	c.clearTokenCookie(rw)

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ConfirmPasswordChange(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ConfirmPasswordChangeParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ConfirmPasswordChangeResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_ConfirmPasswordChange, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	// Clear HTTP cookies.
	c.clearTokenCookie(rw)

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ConfirmRegistration(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ConfirmRegistrationParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ConfirmRegistrationResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_ConfirmRegistration, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) GetSelfRoles(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.GetSelfRolesParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.GetSelfRolesResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_GetSelfRoles, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) GetUserName(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.GetUserNameParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.GetUserNameResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_GetUserName, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) GetUserParameters(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.GetUserParametersParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.GetUserParametersResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_GetUserParameters, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) GetUserRoles(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.GetUserRolesParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.GetUserRolesResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_GetUserRoles, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) GetUserSession(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.GetUserSessionParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.GetUserSessionResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_GetUserSession, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) IsUserLoggedIn(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.IsUserLoggedInParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.IsUserLoggedInResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_IsUserLoggedIn, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ListRegistrationRequestsRFA(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ListRegistrationRequestsRFAParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ListRegistrationRequestsRFAResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_ListRegistrationRequestsRFA, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ListUsers(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ListUsersParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ListUsersResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_ListUsers, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ListUserSessions(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ListUserSessionsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ListUserSessionsResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_ListUserSessions, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) LogUserOutA(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.LogUserOutAParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.LogUserOutAResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_LogUserOutA, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) RejectRegistrationRequestRFA(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.RejectRegistrationRequestRFAParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.RejectRegistrationRequestRFAResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_RejectRegistrationRequestRFA, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) SetUserRoleAuthor(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.SetUserRoleAuthorParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.SetUserRoleAuthorResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_SetUserRoleAuthor, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) SetUserRoleReader(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.SetUserRoleReaderParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.SetUserRoleReaderResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_SetUserRoleReader, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) SetUserRoleWriter(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.SetUserRoleWriterParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.SetUserRoleWriterResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_SetUserRoleWriter, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) StartEmailChange(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.StartEmailChangeParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.StartEmailChangeResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_StartEmailChange, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) StartLogIn(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.StartLogInParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.StartLogInResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_StartLogIn, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) StartLogOut(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.StartLogOutParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.StartLogOutResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_StartLogOut, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) StartPasswordChange(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.StartPasswordChangeParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.StartPasswordChangeResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_StartPasswordChange, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) StartRegistration(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.StartRegistrationParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.StartRegistrationResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_StartRegistration, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) UnbanUser(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.UnbanUserParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.UnbanUserResult)
	var re *jrm1.RpcError
	re, err = c.far.authServiceClient.MakeRequest(context.Background(), rm.Func_UnbanUser, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}

// MessageService.

func (c *Controller) AddForum(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.AddForumParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.AddForumResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_AddForum, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) AddMessage(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.AddMessageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.AddMessageResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_AddMessage, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) AddThread(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.AddThreadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.AddThreadResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_AddThread, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ChangeForumName(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ChangeForumNameParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ChangeForumNameResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_ChangeForumName, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ChangeMessageText(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ChangeMessageTextParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ChangeMessageTextResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_ChangeMessageText, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ChangeMessageThread(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ChangeMessageThreadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ChangeMessageThreadResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_ChangeMessageThread, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ChangeThreadForum(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ChangeThreadForumParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ChangeThreadForumResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_ChangeThreadForum, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ChangeThreadName(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ChangeThreadNameParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ChangeThreadNameResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_ChangeThreadName, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) DeleteForum(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.DeleteForumParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.DeleteForumResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_DeleteForum, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) DeleteMessage(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.DeleteMessageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.DeleteMessageResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_DeleteMessage, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) DeleteThread(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.DeleteThreadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.DeleteThreadResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_DeleteThread, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) GetForum(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.GetForumParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.GetForumResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_GetForum, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) GetMessage(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.GetMessageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.GetMessageResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_GetMessage, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) GetThread(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.GetThreadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.GetThreadResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_GetThread, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ListForums(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ListForumsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ListForumsResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_ListForums, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ListMessages(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ListMessagesParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ListMessagesResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_ListMessages, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) ListThreads(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.ListThreadsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.ListThreadsResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_ListThreads, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) MoveForumDown(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.MoveForumDownParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.MoveForumDownResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_MoveForumDown, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
func (c *Controller) MoveForumUp(ar *rm.Request, _ *http.Request, rw http.ResponseWriter) {
	var err error
	var params rm.MoveForumUpParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}

	params.CommonParams = rm.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(rm.MoveForumUpResult)
	var re *jrm1.RpcError
	re, err = c.far.messageServiceClient.MakeRequest(context.Background(), rm.Func_MoveForumUp, params, result)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}
	if re != nil {
		c.processRpcError(re, rw)
		return
	}

	result.CommonResult.Clear()
	var response = &rm.Response{
		Action: ar.Action,
		Result: result,
	}
	c.respondWithJsonObject(rw, response)
	return
}
