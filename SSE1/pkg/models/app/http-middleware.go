package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	jwtHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/jwt"
	netHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/net"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/auth"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/http/request"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/session"
)

// Context Keys.
const (
	ContextKeyAuthData = "AuthData"
)

// Checks the HTTP Protocol Settings.
func (app *Application) httpProtocolCheck(
	h httprouter.Handle,
) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		ps httprouter.Params,
	) {
		var err = app.checkHttpProtocol(r)
		if err != nil {
			app.handleBadRequestError(w, err, SenderHttpProtocolCheck)
			return
		}
		h(w, r, ps)
		return
	}
}

// Checks the Authentication Data.
func (app *Application) httpAuthentication(
	h httprouter.Handle,
	isAuthRequired bool,
) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		ps httprouter.Params,
	) {
		if !isAuthRequired {
			h(w, r, ps)
			return
		}

		// Get the Token Data and an active Session.
		// If a Session does not exist or is expired, we fail here.
		var ses *session.Session
		var td *jwtHelper.TokenData
		var err error
		var accessIsForbidden bool
		ses, td, accessIsForbidden, err = app.getSessionAndToken(r)
		if err != nil {
			if accessIsForbidden {
				err = errors.New(ErrAuthenticationFailure)
				app.handleForbiddenError(w, err, SenderHttpAuthentication)
			} else {
				app.handleCriticalError(w, err, SenderHttpAuthentication)
			}
			return
		}
		if accessIsForbidden {
			err = errors.New(ErrAuthenticationFailure)
			app.handleForbiddenError(w, err, SenderHttpAuthentication)
			return
		}

		// Read the Client's Hostname and Browser's 'User Agent' Field.
		var host = netHelper.GetAddressHost(r.RemoteAddr)
		err = request.ValidateMachineHost(host)
		if err != nil {
			err = errors.New(ErrAuthenticationFailure)
			app.handleForbiddenError(w, err, SenderHttpAuthentication)
			return
		}
		var bua = r.UserAgent()
		err = request.ValidateMachineBrowserUserAgent(bua)
		if err != nil {
			err = errors.New(ErrAuthenticationFailure)
			app.handleForbiddenError(w, err, SenderHttpAuthentication)
			return
		}
		var buaId uint
		buaId, err = app.buam.GetBrowserUserAgentId(bua)
		if err != nil {
			err = errors.New(ErrAuthenticationFailure)
			app.handleCriticalError(w, err, SenderHttpAuthentication)
			return
		}

		// Verify the Hostname and BUA with the Session.
		if (ses.User.Host != host) ||
			(ses.User.BrowserUserAgentId != buaId) {
			err = errors.New(ErrAuthenticationFailure)
			app.handleForbiddenError(w, err, SenderHttpAuthentication)
			return
		}

		// Pass the following Data to other Handlers in Chain:
		//	- Session;
		//	- Token;
		//	- Hostname, Browser's U.A.
		var ctx = context.WithValue(
			context.Background(),
			ContextKeyAuthData,
			auth.AuthData{
				Session:   ses,
				TokenData: td,
				Machine: &request.UserLogRequestMachine{
					Host: host,
					BrowserUserAgent: request.UserLogRequestMachineBrowserUserAgent{
						Id:   buaId,
						Name: bua,
					},
				},
			},
		)
		r = r.WithContext(ctx)
		h(w, r, ps)
		return
	}
}
