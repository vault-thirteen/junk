package server

import (
	"encoding/json"
	"errors"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/rpc"
	api2 "github.com/vault-thirteen/SimpleBB/pkg/GWM/api"
	ch "github.com/vault-thirteen/SimpleBB/pkg/common/models/http"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"io"
	"log"
	"net/http"

	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	s "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
	"github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

const (
	// ACM.
	ApiFunctionName_RegisterUser                           = "registerUser"
	ApiFunctionName_GetListOfRegistrationsReadyForApproval = "getListOfRegistrationsReadyForApproval"
	ApiFunctionName_RejectRegistrationRequest              = "rejectRegistrationRequest"
	ApiFunctionName_ApproveAndRegisterUser                 = "approveAndRegisterUser"
	ApiFunctionName_LogUserIn                              = "logUserIn"
	ApiFunctionName_LogUserOut                             = "logUserOut"
	ApiFunctionName_LogUserOutA                            = "logUserOutA"
	ApiFunctionName_GetListOfLoggedUsers                   = "getListOfLoggedUsers"
	ApiFunctionName_GetListOfLoggedUsersOnPage             = "getListOfLoggedUsersOnPage"
	ApiFunctionName_GetListOfAllUsers                      = "getListOfAllUsers"
	ApiFunctionName_GetListOfAllUsersOnPage                = "getListOfAllUsersOnPage"
	ApiFunctionName_IsUserLoggedIn                         = "isUserLoggedIn"
	ApiFunctionName_ChangePassword                         = "changePassword"
	ApiFunctionName_ChangeEmail                            = "changeEmail"
	ApiFunctionName_GetUserSession                         = "getUserSession"
	ApiFunctionName_GetUserName                            = "getUserName"
	ApiFunctionName_GetUserRoles                           = "getUserRoles"
	ApiFunctionName_ViewUserParameters                     = "viewUserParameters"
	ApiFunctionName_SetUserRoleAuthor                      = "setUserRoleAuthor"
	ApiFunctionName_SetUserRoleWriter                      = "setUserRoleWriter"
	ApiFunctionName_SetUserRoleReader                      = "setUserRoleReader"
	ApiFunctionName_GetSelfRoles                           = "getSelfRoles"
	ApiFunctionName_BanUser                                = "banUser"
	ApiFunctionName_UnbanUser                              = "unbanUser"

	// MM.
	ApiFunctionName_AddSection                  = "addSection"
	ApiFunctionName_ChangeSectionName           = "changeSectionName"
	ApiFunctionName_ChangeSectionParent         = "changeSectionParent"
	ApiFunctionName_GetSection                  = "getSection"
	ApiFunctionName_MoveSectionUp               = "moveSectionUp"
	ApiFunctionName_MoveSectionDown             = "moveSectionDown"
	ApiFunctionName_DeleteSection               = "deleteSection"
	ApiFunctionName_AddForum                    = "addForum"
	ApiFunctionName_ChangeForumName             = "changeForumName"
	ApiFunctionName_ChangeForumSection          = "changeForumSection"
	ApiFunctionName_GetForum                    = "getForum"
	ApiFunctionName_MoveForumUp                 = "moveForumUp"
	ApiFunctionName_MoveForumDown               = "moveForumDown"
	ApiFunctionName_DeleteForum                 = "deleteForum"
	ApiFunctionName_AddThread                   = "addThread"
	ApiFunctionName_ChangeThreadName            = "changeThreadName"
	ApiFunctionName_ChangeThreadForum           = "changeThreadForum"
	ApiFunctionName_GetThread                   = "getThread"
	ApiFunctionName_GetThreadNamesByIds         = "getThreadNamesByIds"
	ApiFunctionName_MoveThreadUp                = "moveThreadUp"
	ApiFunctionName_MoveThreadDown              = "moveThreadDown"
	ApiFunctionName_DeleteThread                = "deleteThread"
	ApiFunctionName_AddMessage                  = "addMessage"
	ApiFunctionName_ChangeMessageText           = "changeMessageText"
	ApiFunctionName_ChangeMessageThread         = "changeMessageThread"
	ApiFunctionName_GetMessage                  = "getMessage"
	ApiFunctionName_GetLatestMessageOfThread    = "getLatestMessageOfThread"
	ApiFunctionName_DeleteMessage               = "deleteMessage"
	ApiFunctionName_ListThreadAndMessages       = "listThreadAndMessages"
	ApiFunctionName_ListThreadAndMessagesOnPage = "listThreadAndMessagesOnPage"
	ApiFunctionName_ListForumAndThreads         = "listForumAndThreads"
	ApiFunctionName_ListForumAndThreadsOnPage   = "listForumAndThreadsOnPage"
	ApiFunctionName_ListSectionsAndForums       = "listSectionsAndForums"

	// NM.
	ApiFunctionName_AddNotification             = "addNotification"
	ApiFunctionName_GetNotification             = "getNotification"
	ApiFunctionName_GetNotifications            = "getNotifications"
	ApiFunctionName_GetNotificationsOnPage      = "getNotificationsOnPage"
	ApiFunctionName_GetUnreadNotifications      = "getUnreadNotifications"
	ApiFunctionName_CountUnreadNotifications    = "countUnreadNotifications"
	ApiFunctionName_MarkNotificationAsRead      = "markNotificationAsRead"
	ApiFunctionName_DeleteNotification          = "deleteNotification"
	ApiFunctionName_AddResource                 = "addResource"
	ApiFunctionName_GetResource                 = "getResource"
	ApiFunctionName_GetResourceValue            = "getResourceValue"
	ApiFunctionName_GetListOfAllResourcesOnPage = "getListOfAllResourcesOnPage"
	ApiFunctionName_DeleteResource              = "deleteResource"

	// SM.
	ApiFunctionName_AddSubscription            = "addSubscription"
	ApiFunctionName_IsSelfSubscribed           = "isSelfSubscribed"
	ApiFunctionName_IsUserSubscribed           = "isUserSubscribed"
	ApiFunctionName_CountSelfSubscriptions     = "countSelfSubscriptions"
	ApiFunctionName_GetSelfSubscriptions       = "getSelfSubscriptions"
	ApiFunctionName_GetSelfSubscriptionsOnPage = "getSelfSubscriptionsOnPage"
	ApiFunctionName_GetUserSubscriptions       = "getUserSubscriptions"
	ApiFunctionName_GetUserSubscriptionsOnPage = "getUserSubscriptionsOnPage"
	ApiFunctionName_DeleteSelfSubscription     = "deleteSelfSubscription"
	ApiFunctionName_DeleteSubscription         = "deleteSubscription"
)

func (srv *Server) handlePublicSettings(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		srv.respondMethodNotAllowed(rw)
		return
	}

	// Check accepted MIME types.
	ok, err := hh.CheckBrowserSupportForJson(req)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}
	if !ok {
		srv.respondNotAcceptable(rw)
		return
	}

	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.Header().Set(header.HttpHeaderContentType, ch.ContentType_Json)

	_, err = rw.Write(srv.publicSettingsFileData)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (srv *Server) handlePublicApi(rw http.ResponseWriter, req *http.Request, clientIPA simple.IPAS) {
	if req.Method != http.MethodPost {
		srv.respondMethodNotAllowed(rw)
		return
	}

	// Check accepted MIME types.
	ok, err := hh.CheckBrowserSupportForJson(req)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}

	if !ok {
		srv.respondNotAcceptable(rw)
		return
	}

	var reqBody []byte
	reqBody, err = io.ReadAll(req.Body)
	if err != nil {
		srv.processInternalServerError(rw, err)
		return
	}

	// Check the action.
	var arwoa api2.RequestWithOnlyAction
	err = json.Unmarshal(reqBody, &arwoa)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}

	if (arwoa.Action == nil) ||
		(arwoa.Parameters == nil) {
		srv.respondBadRequest(rw)
		return
	}

	var handler api2.RequestHandler
	handler, ok = srv.apiHandlers[*arwoa.Action]
	if !ok {
		srv.respondNotFound(rw)
		return
	}

	var token *simple.WebTokenString
	token, err = simple.GetToken(req)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}

	var ar = &api2.Request{
		Action:     arwoa.Action,
		Parameters: arwoa.Parameters,
		Authorisation: &cmr.Auth{
			UserIPA: clientIPA,
		},
	}

	if token != nil {
		ar.Authorisation.Token = *token
	}

	handler(ar, req, rw)
	return
}

func (srv *Server) handleCaptcha(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		srv.respondMethodNotAllowed(rw)
		return
	}

	srv.rcsProxy.ServeHTTP(rw, req)
}

// handleFrontEnd serves static files for the front end part.
func (srv *Server) handleFrontEnd(rw http.ResponseWriter, req *http.Request, clientIPA simple.IPAS) {
	// While the number of cases is less than 10..20, the "switch" branching
	// works faster than other methods.
	switch simple.Path(req.URL.Path) {
	case srv.frontEnd.AdminHtmlPage.UrlPath:
		srv.handleAdminFrontEnd(rw, req, clientIPA)
		return

	case srv.frontEnd.AdminJs.UrlPath:
		srv.handleAdminFrontEnd(rw, req, clientIPA)
		return

	case srv.frontEnd.ApiJs.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.ApiJs)
		return

	case srv.frontEnd.ArgonJs.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.ArgonJs)
		return

	case srv.frontEnd.ArgonWasm.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.ArgonWasm)
		return

	case srv.frontEnd.BppJs.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.BppJs)
		return

	case srv.frontEnd.IndexHtmlPage.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.IndexHtmlPage)
		return

	case srv.frontEnd.LoaderScript.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.LoaderScript)
		return

	case srv.frontEnd.CssStyles.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.CssStyles)
		return

	default:
		srv.respondNotFound(rw)
		return
	}
}

// handleAdminFrontEnd serves static files for the front end part for
// administrator users.
func (srv *Server) handleAdminFrontEnd(rw http.ResponseWriter, req *http.Request, clientIPA simple.IPAS) {
	// Step 1. Filter for fake URL.
	switch simple.Path(req.URL.Path) {
	case s.FrontEndAdminPath: // <- /admin
	case srv.frontEnd.AdminHtmlPage.UrlPath: // <- /fe/admin.html
	case srv.frontEnd.AdminJs.UrlPath: // <- /fe/admin.js
	default:
		srv.respondNotFound(rw)
		return
	}

	// Step 2. Filter for administrator role.
	ok, err := srv.checkForAdministrator(rw, req, clientIPA)
	if err != nil {
		srv.processInternalServerError(rw, err)
		return
	}
	if !ok {
		srv.respondForbidden(rw)
		return
	}

	// Step 3. Page selection.
	switch simple.Path(req.URL.Path) {
	case s.FrontEndAdminPath: // <- /admin
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.AdminHtmlPage)
		return

	case srv.frontEnd.AdminHtmlPage.UrlPath: // <- /fe/admin.html
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.AdminHtmlPage)
		return

	case srv.frontEnd.AdminJs.UrlPath: // <- /fe/admin.js
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.AdminJs)
		return

	default:
		srv.respondNotFound(rw)
		return
	}
}

func (srv *Server) handleFrontEndStaticFile(rw http.ResponseWriter, req *http.Request, fedf models.FrontEndFileData) {
	if req.Method != http.MethodGet {
		srv.respondMethodNotAllowed(rw)
		return
	}

	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.Header().Set(header.HttpHeaderContentType, fedf.ContentType)

	_, err := rw.Write(fedf.CachedFile)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (srv *Server) checkForAdministrator(rw http.ResponseWriter, req *http.Request, clientIPA simple.IPAS) (ok bool, err error) {
	var token *simple.WebTokenString
	token, err = simple.GetToken(req)
	if err != nil {
		srv.respondBadRequest(rw)
		return false, err
	}

	if token == nil {
		srv.respondForbidden(rw)
		return
	}

	// Make a 'GetSelfRoles' RPC request to verify user's roles.
	action := ApiFunctionName_GetSelfRoles
	var params json.RawMessage = []byte("{}")
	var ar = &api2.Request{
		Action:     &action,
		Parameters: &params,
		Authorisation: &cmr.Auth{
			UserIPA: clientIPA,
			Token:   *token,
		},
	}

	var response *api2.Response
	response, err = srv.getSelfRoles(ar)
	if err != nil {
		srv.processInternalServerError(rw, err)
		return false, err
	}

	var result *am.GetSelfRolesResult
	result, ok = response.Result.(*am.GetSelfRolesResult)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(rw, err)
		return false, err
	}

	if !result.User.GetUserParameters().GetRoles().IsAdministrator {
		srv.respondForbidden(rw)
		return false, nil
	}

	return true, nil
}
