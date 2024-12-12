package server

import (
	"context"
	"encoding/json"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/rpc"
	api2 "github.com/vault-thirteen/SimpleBB/pkg/GWM/api"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/rpc"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/rpc"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	mc "github.com/vault-thirteen/SimpleBB/pkg/MM/client"
	nc "github.com/vault-thirteen/SimpleBB/pkg/NM/client"
	sc "github.com/vault-thirteen/SimpleBB/pkg/SM/client"
)

// Unfortunately, Go language still has very poor support for generic
// programming. The D.R.Y. principle is, thus, violated in this file.
// All the functions in this file are created by a copy-paste method with some
// minor exceptions. The only exceptions are listed below.
//
//	1. The 'LogUserIn' function has additional code which:
// 		1.1. ignores a token;
//		1.2. sets a token.
//
//	2. The 'LogUserOut' function has additional code which:
//		2.1. clears a token.
//
//	3. The 'ChangePassword' function has additional code which:
//		3.1. clears a token.
//
//	4. The 'ChangeEmail' function has additional code which:
//		4.1. clears a token.

// Service functions.

// ACM.

func (srv *Server) RegisterUser(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.RegisterUserParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.RegisterUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncRegisterUser, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetListOfRegistrationsReadyForApproval(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.GetListOfRegistrationsReadyForApprovalParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetListOfRegistrationsReadyForApprovalResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetListOfRegistrationsReadyForApproval, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) RejectRegistrationRequest(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.RejectRegistrationRequestParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.RejectRegistrationRequestResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncRejectRegistrationRequest, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ApproveAndRegisterUser(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.ApproveAndRegisterUserParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.ApproveAndRegisterUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncApproveAndRegisterUser, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) LogUserIn(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.LogUserInParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	// To allow a user with an outdated token to get into the system, we ignore
	// the old token. [1.1]
	params.CommonParams.Auth.Token = ""

	var result = new(am.LogUserInResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncLogUserIn, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}

	// Save the token in HTTP cookies. [1.2]
	if result.IsWebTokenSet {
		srv.setTokenCookie(hrw, result.WebTokenString)
		result.WebTokenString = ""
	}

	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) LogUserOut(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.LogUserOutParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.LogUserOutResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncLogUserOut, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}

	// Clear the token in HTTP cookies. [2.1]
	srv.clearTokenCookie(hrw)

	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) LogUserOutA(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.LogUserOutAParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.LogUserOutAResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncLogUserOutA, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}

	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetListOfLoggedUsers(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.GetListOfLoggedUsersParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetListOfLoggedUsersResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetListOfLoggedUsers, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetListOfLoggedUsersOnPage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.GetListOfLoggedUsersOnPageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetListOfLoggedUsersOnPageResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetListOfLoggedUsersOnPage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetListOfAllUsers(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.GetListOfAllUsersParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetListOfAllUsersResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetListOfAllUsers, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetListOfAllUsersOnPage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.GetListOfAllUsersOnPageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetListOfAllUsersOnPageResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetListOfAllUsersOnPage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) IsUserLoggedIn(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.IsUserLoggedInParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.IsUserLoggedInResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncIsUserLoggedIn, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangePassword(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.ChangePasswordParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.ChangePasswordResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncChangePassword, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}

	// Clear the token in HTTP cookies. [3.1]
	if result.OK {
		srv.clearTokenCookie(hrw)
	}

	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeEmail(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.ChangeEmailParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.ChangeEmailResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncChangeEmail, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}

	// Clear the token in HTTP cookies. [4.1]
	if result.OK {
		srv.clearTokenCookie(hrw)
	}

	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetUserSession(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.GetUserSessionParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetUserSessionResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetUserSession, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}

	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetUserName(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.GetUserNameParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetUserNameResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetUserName, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetUserRoles(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.GetUserRolesParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetUserRolesResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetUserRoles, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ViewUserParameters(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.ViewUserParametersParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.ViewUserParametersResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncViewUserParameters, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) SetUserRoleAuthor(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.SetUserRoleAuthorParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.SetUserRoleAuthorResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncSetUserRoleAuthor, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) SetUserRoleWriter(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.SetUserRoleWriterParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.SetUserRoleWriterResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncSetUserRoleWriter, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) SetUserRoleReader(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.SetUserRoleReaderParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.SetUserRoleReaderResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncSetUserRoleReader, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

// GetSelfRoles is a normal version of 'GetSelfRoles' RPC request for public
// usage. For internal purposes, use its internal variant – 'getSelfRoles'.
func (srv *Server) GetSelfRoles(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.GetSelfRolesParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetSelfRolesResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetSelfRoles, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

// getSelfRoles is a clone of the 'GetSelfRoles' method, intended for internal
// usage only.
func (srv *Server) getSelfRoles(ar *api2.Request) (response *api2.Response, err error) {
	var params am.GetSelfRolesParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		return nil, err
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.GetSelfRolesResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetSelfRoles, params, result)
	if err != nil {
		return nil, err
	}
	if re != nil {
		return nil, re.AsError()
	}

	result.CommonResult.Clear()
	response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	return response, nil
}

func (srv *Server) BanUser(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.BanUserParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.BanUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncBanUser, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) UnbanUser(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params am.UnbanUserParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{
		Auth: ar.Authorisation,
	}

	var result = new(am.UnbanUserResult)
	var re *jrm1.RpcError
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncUnbanUser, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_ACM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{
		Action: ar.Action,
		Result: result,
	}
	srv.respondWithJsonObject(hrw, response)
	return
}

// MM.

func (srv *Server) AddSection(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.AddSectionParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.AddSectionResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncAddSection, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeSectionName(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ChangeSectionNameParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ChangeSectionNameResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncChangeSectionName, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeSectionParent(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ChangeSectionParentParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ChangeSectionParentResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncChangeSectionParent, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetSection(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.GetSectionParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.GetSectionResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncGetSection, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) MoveSectionUp(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.MoveSectionUpParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.MoveSectionUpResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncMoveSectionUp, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) MoveSectionDown(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.MoveSectionDownParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.MoveSectionDownResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncMoveSectionDown, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) DeleteSection(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.DeleteSectionParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.DeleteSectionResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncDeleteSection, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) AddForum(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.AddForumParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.AddForumResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncAddForum, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeForumName(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ChangeForumNameParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ChangeForumNameResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncChangeForumName, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeForumSection(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ChangeForumSectionParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ChangeForumSectionResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncChangeForumSection, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetForum(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.GetForumParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.GetForumResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncGetForum, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) MoveForumUp(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.MoveForumUpParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.MoveForumUpResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncMoveForumUp, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) MoveForumDown(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.MoveForumDownParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.MoveForumDownResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncMoveForumDown, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) DeleteForum(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.DeleteForumParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.DeleteForumResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncDeleteForum, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) AddThread(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.AddThreadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.AddThreadResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncAddThread, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeThreadName(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ChangeThreadNameParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ChangeThreadNameResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncChangeThreadName, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeThreadForum(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ChangeThreadForumParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ChangeThreadForumResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncChangeThreadForum, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetThread(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.GetThreadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.GetThreadResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncGetThread, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetThreadNamesByIds(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.GetThreadNamesByIdsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.GetThreadNamesByIdsResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncGetThreadNamesByIds, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) MoveThreadUp(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.MoveThreadUpParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.MoveThreadUpResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncMoveThreadUp, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) MoveThreadDown(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.MoveThreadDownParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.MoveThreadDownResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncMoveThreadDown, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) DeleteThread(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.DeleteThreadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.DeleteThreadResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncDeleteThread, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) AddMessage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.AddMessageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.AddMessageResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncAddMessage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeMessageText(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ChangeMessageTextParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ChangeMessageTextResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncChangeMessageText, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ChangeMessageThread(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ChangeMessageThreadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ChangeMessageThreadResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncChangeMessageThread, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetMessage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.GetMessageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.GetMessageResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncGetMessage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetLatestMessageOfThread(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.GetLatestMessageOfThreadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.GetLatestMessageOfThreadResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncGetLatestMessageOfThread, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) DeleteMessage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.DeleteMessageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.DeleteMessageResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncDeleteMessage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ListThreadAndMessages(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ListThreadAndMessagesParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ListThreadAndMessagesResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncListThreadAndMessages, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ListThreadAndMessagesOnPage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ListThreadAndMessagesOnPageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ListThreadAndMessagesOnPageResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncListThreadAndMessagesOnPage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ListForumAndThreads(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ListForumAndThreadsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ListForumAndThreadsResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncListForumAndThreads, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ListForumAndThreadsOnPage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ListForumAndThreadsOnPageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ListForumAndThreadsOnPageResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncListForumAndThreadsOnPage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) ListSectionsAndForums(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params mm.ListSectionsAndForumsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(mm.ListSectionsAndForumsResult)
	var re *jrm1.RpcError
	re, err = srv.mmServiceClient.MakeRequest(context.Background(), mc.FuncListSectionsAndForums, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_MM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

// NM.

func (srv *Server) AddNotification(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.AddNotificationParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.AddNotificationResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncAddNotification, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetNotification(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.GetNotificationParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.GetNotificationResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncGetNotification, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetNotifications(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.GetNotificationsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.GetNotificationsResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncGetNotifications, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetNotificationsOnPage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.GetNotificationsOnPageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.GetNotificationsOnPageResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncGetNotificationsOnPage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetUnreadNotifications(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.GetUnreadNotificationsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.GetUnreadNotificationsResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncGetUnreadNotifications, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) CountUnreadNotifications(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.CountUnreadNotificationsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.CountUnreadNotificationsResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncCountUnreadNotifications, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) MarkNotificationAsRead(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.MarkNotificationAsReadParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.MarkNotificationAsReadResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncMarkNotificationAsRead, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) DeleteNotification(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.DeleteNotificationParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.DeleteNotificationResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncDeleteNotification, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) AddResource(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.AddResourceParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.AddResourceResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncAddResource, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetResource(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.GetResourceParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.GetResourceResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncGetResource, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetResourceValue(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.GetResourceValueParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.GetResourceValueResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncGetResourceValue, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetListOfAllResourcesOnPage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.GetListOfAllResourcesOnPageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.GetListOfAllResourcesOnPageResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncGetListOfAllResourcesOnPage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) DeleteResource(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params nm.DeleteResourceParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(nm.DeleteResourceResult)
	var re *jrm1.RpcError
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncDeleteResource, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_NM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

// SM.

func (srv *Server) AddSubscription(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.AddSubscriptionParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.AddSubscriptionResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncAddSubscription, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) IsSelfSubscribed(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.IsSelfSubscribedParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.IsSelfSubscribedResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncIsSelfSubscribed, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) IsUserSubscribed(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.IsUserSubscribedParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.IsUserSubscribedResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncIsUserSubscribed, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) CountSelfSubscriptions(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.CountSelfSubscriptionsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.CountSelfSubscriptionsResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncCountSelfSubscriptions, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetSelfSubscriptions(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.GetSelfSubscriptionsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.GetSelfSubscriptionsResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncGetSelfSubscriptions, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetSelfSubscriptionsOnPage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.GetSelfSubscriptionsOnPageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.GetSelfSubscriptionsOnPageResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncGetSelfSubscriptionsOnPage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetUserSubscriptions(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.GetUserSubscriptionsParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.GetUserSubscriptionsResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncGetUserSubscriptions, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) GetUserSubscriptionsOnPage(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.GetUserSubscriptionsOnPageParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.GetUserSubscriptionsOnPageResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncGetUserSubscriptionsOnPage, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) DeleteSelfSubscription(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.DeleteSelfSubscriptionParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.DeleteSelfSubscriptionResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncDeleteSelfSubscription, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}

func (srv *Server) DeleteSubscription(ar *api2.Request, _ *http.Request, hrw http.ResponseWriter) {
	var err error
	var params sm.DeleteSubscriptionParams
	err = json.Unmarshal(*ar.Parameters, &params)
	if err != nil {
		srv.respondBadRequest(hrw)
		return
	}

	params.CommonParams = cmr.CommonParams{Auth: ar.Authorisation}

	var result = new(sm.DeleteSubscriptionResult)
	var re *jrm1.RpcError
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncDeleteSubscription, params, result)
	if err != nil {
		srv.processInternalServerError(hrw, err)
		return
	}
	if re != nil {
		srv.processRpcError(app.ModuleId_SM, re, hrw)
		return
	}

	result.CommonResult.Clear()
	var response = &api2.Response{Action: ar.Action, Result: result}
	srv.respondWithJsonObject(hrw, response)
	return
}
