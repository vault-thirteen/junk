package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/http"
	"github.com/vault-thirteen/junk/SSE1/pkg/helper/jwt"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/auth"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/http/request"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/http/response"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/session"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/user"
)

// An HTTP Handler which registers a new User.
func (app *Application) httpHandlerApiUserRegister(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	// No Authentication is required.

	// Read the Request.
	var requestObject *request.UserRegistrationRequest
	var ok bool
	requestObject, ok = app.getHttpRequest_ApiUserRegister(w, r, ps)
	if !ok {
		return
	}

	// Preparation.
	// Check if the User Name is free.
	var nameIsFree bool
	var err error
	nameIsFree, err = app.storage.IsUserAuthenticationNameFree(
		requestObject.User.InternalName,
	)
	if err != nil {
		app.handleCriticalError(w, err, SenderIsUserAuthenticationNameFree)
		return
	}
	if !nameIsFree {
		err = fmt.Errorf(ErrfUserAuthenticationNameIsTaken, requestObject.User.InternalName)
		app.handleForbiddenError(w, err, SenderIsUserAuthenticationNameFree)
		return
	}
	var usr = &user.User{
		Authentication: user.UserAuthentication{
			Name:     requestObject.User.InternalName,
			Password: requestObject.User.Password,
		},
		Registration: user.UserRegistration{
			SecretCode: requestObject.User.SecretCode,
		},
		PublicName: requestObject.User.PublicName,
	}

	// Main Work.
	err = app.storage.RegisterUser(usr)
	if err != nil {
		app.handleCriticalError(w, err, SenderRegisterUser)
		return
	}
}

// An HTTP Handler which disables a registered User.
// This Handler may not be used too frequently.
func (app *Application) httpHandlerApiUserDisable(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	// Unpack the Authentication Data from Context.
	var authData auth.AuthData
	var err error
	authData, err = unpackAuthDataFromContext(r)
	if err != nil {
		app.handleCriticalError(w, err, SenderUnpackAuthDataFromContext)
		return
	}

	// Read the Request.
	var requestObject *request.UserDisablingRequest
	var ok bool
	requestObject, ok = app.getHttpRequest_ApiUserDisable(w, r, ps)
	if !ok {
		return
	}

	// Preparation.
	var usr = &user.User{
		//Id:// See below.
		Authentication: user.UserAuthentication{
			Name:     requestObject.User.InternalName,
			Password: requestObject.User.Password,
		},
		Registration: user.UserRegistration{
			SecretCode: requestObject.User.SecretCode,
		},
	}
	usr.Id, err = app.storage.GetUserIdByAuthenticationName(usr.Authentication.Name)
	if err != nil {
		app.handleCriticalError(w, err, SenderGetUserIdByAuthenticationName)
		return
	}

	// We can only disable ourselves.
	if usr.Id != authData.Session.User.Id {
		err = errors.New(ErrUserCanNotBeDisabled)
		app.handleForbiddenError(w, err, SenderHttpHandlerApiUserDisable)
		return
	}

	// Check if a registered User with Id exists.
	ok, err = app.storage.RegisteredUserIdExists(usr.Id)
	if err != nil {
		app.handleCriticalError(w, err, SenderRegisteredUserIdExists)
		return
	}
	if !ok {
		err = errors.New(ErrRegisteredUserDoesNotExist)
		app.handleForbiddenError(w, err, SenderRegisteredUserIdExists)
		return
	}

	// Try to disable a User.
	// The Process fails if you try to do it too frequently.
	err = app.storage.DisableUser(usr)
	if err != nil {
		app.handleForbiddenError(w, err, SenderDisableUser)
		return
	}
}

// An HTTP Handler which logs the registered User into the Server.
// This Handler may not be used too frequently.
func (app *Application) httpHandlerApiUserLogIn(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	// No Authentication is required.

	// Read the Request.
	var requestObject *request.UserLogInRequest
	var ok bool
	requestObject, ok = app.getHttpRequest_ApiUserLogIn(w, r, ps)
	if !ok {
		return
	}

	// Preparation.
	var usr = &user.User{
		//Id: // See below.
		Authentication: user.UserAuthentication{
			Name:     requestObject.User.InternalName,
			Password: requestObject.User.Password,
		},
	}
	var err error
	usr.Id, err = app.storage.GetUserIdByAuthenticationName(usr.Authentication.Name)
	if err != nil {
		app.handleForbiddenError(w, err, SenderGetUserIdByAuthenticationName)
		return
	}

	// Check whether a registered User exists.
	ok, err = app.storage.RegisteredUserIdExists(usr.Id)
	if err != nil {
		app.handleCriticalError(w, err, SenderRegisteredUserIdExists)
		return
	}
	if !ok {
		err = errors.New(ErrRegisteredUserDoesNotExist)
		app.handleForbiddenError(w, err, SenderRegisteredUserIdExists)
		return
	}

	// Try to log the User in.
	// The Process fails if you try to do it too frequently.
	var ses *session.Session
	var token *jwt.TokenData
	ses, token, err = app.storage.LogUserIn(usr, &requestObject.Machine)
	if err != nil {
		app.handleForbiddenError(w, err, SenderLogUserIn)
		return
	}

	// Respond.
	app.setTokenHeader(w, token.TokenString)
	setMarkerCookie(w, ses.Marker, app.configuration.Server.HttpServer.CookiePath)
}

// An HTTP Handler which logs the registered User out of the Server.
func (app *Application) httpHandlerApiUserLogOut(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	// Unpack the Authentication Data from Context.
	var authData auth.AuthData
	var err error
	authData, err = unpackAuthDataFromContext(r)
	if err != nil {
		app.handleCriticalError(w, err, SenderUnpackAuthDataFromContext)
		return
	}

	// Read the Request.
	var requestObject *request.UserLogOutRequest
	var ok bool
	requestObject, ok = app.getHttpRequest_ApiUserLogOut(w, r, ps)
	if !ok {
		return
	}

	// Preparation.
	var usr = &user.User{
		//Id:// See below.
		Authentication: user.UserAuthentication{
			Name: requestObject.User.InternalName,
		},
	}
	usr.Id, err = app.storage.GetUserIdByAuthenticationName(usr.Authentication.Name)
	if err != nil {
		app.handleCriticalError(w, err, SenderGetUserIdByAuthenticationName)
		return
	}

	// We can only log out ourselves.
	if usr.Id != authData.Session.User.Id {
		err = errors.New(ErrCanNotLogOutOtherUsers)
		app.handleForbiddenError(w, err, SenderHttpHandlerApiUserLogOut)
	}

	// Check whether a registered User exists.
	ok, err = app.storage.RegisteredUserIdExists(usr.Id)
	if err != nil {
		app.handleCriticalError(w, err, SenderRegisteredUserIdExists)
		return
	}
	if !ok {
		err = errors.New(ErrRegisteredUserDoesNotExist)
		app.handleForbiddenError(w, err, SenderRegisteredUserIdExists)
		return
	}

	// Try to log the User out.
	err = app.storage.LogUserOut(usr, authData.Session)
	if err != nil {
		app.handleForbiddenError(w, err, SenderLogUserOut)
		return
	}
}

// An HTTP Handler which lists public Names of all registered Users.
// Returns an Error for non-authenticated Clients.
func (app *Application) httpHandlerApiUsersList(
	w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params,
) {
	// Unpack the Authentication Data from Context.
	var authData auth.AuthData
	var err error
	authData, err = unpackAuthDataFromContext(r)
	if err != nil {
		app.handleCriticalError(w, err, SenderUnpackAuthDataFromContext)
		return
	}

	// Get the List of Names of registered Users.
	var responseObject = new(response.UsersPublicNameList)
	responseObject.UsersPublicNames, err = app.storage.ListRegisteredUsersPublicNames()
	if err != nil {
		app.handleCriticalError(w, err, SenderListRegisteredUsersPublicNames)
		return
	}

	// Update the Session.
	err = app.storage.UpdateActiveSessionLastAccessTime(authData.Session)
	if err != nil {
		app.handleCriticalError(w, err, SenderUpdateActiveSessionLastAccessTime)
		return
	}

	// Respond.
	err = httpHelper.RespondWithJsonObject(w, responseObject)
	if err != nil {
		app.handleCriticalError(w, err, SenderRespondWithJsonObject)
		return
	}
}
