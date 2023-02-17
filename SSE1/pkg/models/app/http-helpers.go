package app

import (
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	httpHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/http"
	jwtHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/jwt"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/auth"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/session"
)

// Errors.
const (
	ErrTlsVersionTooLow  = "TLS Version is too low: %v vs %v"
	ErrTlsVersionTooHigh = "TLS Version is too high: %v vs %v"
)

// Settings.
const (
	CookieMarker         = "MRK"
	ErrorSenderDelimiter = ": "
)

// Responds with a critical Error ('Internal Server Error').
func (app *Application) handleCriticalError(
	w http.ResponseWriter,
	criticalError error,
	sender string,
) {
	criticalError = errors.New(sender + ErrorSenderDelimiter + criticalError.Error())
	var err = fmt.Errorf(ErrfCriticalError, criticalError)
	app.errorChannel <- err
	var responseError = httpHelper.SendHttpResponseInternalServerError(w, err)
	if responseError != nil {
		var err2 = fmt.Errorf(ErrfHttpResponseError, responseError)
		app.errorChannel <- err2
	}
}

// Responds with a 'Bad Request Error'.
func (app *Application) handleBadRequestError(
	w http.ResponseWriter,
	badRequestError error,
	sender string,
) {
	badRequestError = errors.New(sender + ErrorSenderDelimiter + badRequestError.Error())
	var err = fmt.Errorf(ErrfBadRequestError, badRequestError)
	app.errorChannel <- err
	var responseError = httpHelper.SendHttpResponseBadRequestError(w, err)
	if responseError != nil {
		var err2 = fmt.Errorf(ErrfHttpResponseError, responseError)
		app.errorChannel <- err2
	}
}

// Responds with a 'Forbidden Error'.
func (app *Application) handleForbiddenError(
	w http.ResponseWriter,
	forbiddenError error,
	sender string,
) {
	forbiddenError = errors.New(sender + ErrorSenderDelimiter + forbiddenError.Error())
	var err = fmt.Errorf(ErrfForbiddenError, forbiddenError)
	app.errorChannel <- err
	var responseError = httpHelper.SendHttpResponseForbiddenError(w, err)
	if responseError != nil {
		var err2 = fmt.Errorf(ErrfHttpResponseError, responseError)
		app.errorChannel <- err2
	}
}

// Checks the TLS Version of the Protocol.
func checkTlsProtocolVersion(
	r *http.Request,
) (err error) {
	const (
		TlsVersionMin = tls.VersionTLS13
		TlsVersionMax = tls.VersionTLS13 // This may change in Future.
	)
	if r.TLS.Version < TlsVersionMin {
		err = fmt.Errorf(
			ErrTlsVersionTooLow,
			r.TLS.Version,
			TlsVersionMin,
		)
	}
	if r.TLS.Version > TlsVersionMax {
		err = fmt.Errorf(
			ErrTlsVersionTooHigh,
			r.TLS.Version,
			TlsVersionMax,
		)
	}
	return
}

// Checks the HTTP Protocol Settings.
func (app *Application) checkHttpProtocol(
	r *http.Request,
) (err error) {
	if app.configuration.Server.TLS.IsEnabled {
		err = checkTlsProtocolVersion(r)
		if err != nil {
			return
		}
		return
	}
	return
}

// Sets a Token Header to the HTTP Response.
func (app *Application) setTokenHeader(
	w http.ResponseWriter,
	tokenString string,
) {
	httpHelper.SetHttpResponseHeader(
		w,
		app.configuration.Server.HttpServer.TokenHeader,
		tokenString,
	)
}

// Gets a Token Header from the HTTP Request.
func (app *Application) getTokenHeader(
	r *http.Request,
) (tokenString string, err error) {
	return httpHelper.GetHttpRequestHeader(
		r,
		app.configuration.Server.HttpServer.TokenHeader,
	)
}

// Sets the Marker's Cookie in the HTTP Response.
func setMarkerCookie(
	w http.ResponseWriter,
	marker string,
	cookiePath string,
) {
	httpHelper.SetHttpSecureCookie(w, CookieMarker, marker, cookiePath)
}

// Reads the Marker's Cookie from the HTTP Response.
func getMarkerCookie(
	r *http.Request,
) (marker string, err error) {
	return httpHelper.GetHttpCookie(r, CookieMarker)
}

// Unpacks the Authentication Data from Context.
func unpackAuthDataFromContext(
	r *http.Request,
) (ad auth.AuthData, err error) {
	var authDataIfc = r.Context().Value(ContextKeyAuthData)
	var ok bool
	ad, ok = authDataIfc.(auth.AuthData)
	if !ok {
		err = errors.New(ErrTypeCastFailure)
		return
	}
	return
}

// Gets Information about Token and Session from the HTTP Request and Storage.
// While Token only stores an Id of a Session, the most Part of the Session is
// taken from the Storage. The full Marker is taken from an HTTP Cookie, but
// the Marker's Hash Sum is taken from the Token. When all the required Data is
// received, the Method verifies it and returns an Error on Failure.
func (app *Application) getSessionAndToken(
	r *http.Request,
) (ses *session.Session, td *jwtHelper.TokenData, accessIsForbidden bool, err error) {
	var marker string
	accessIsForbidden = true
	marker, err = getMarkerCookie(r)
	if err != nil {
		return
	}
	var tokenString string
	tokenString, err = app.getTokenHeader(r)
	if err != nil {
		return
	}
	var token *jwt.Token
	token, err = jwtHelper.VerifyJWToken(tokenString, app.getTokenKey)
	if err != nil {
		return
	}
	td = &jwtHelper.TokenData{
		UniqueMarker: marker,
	}
	td.SessionId, td.UniqueMarkerHash, err = jwtHelper.GetTokenData(token)
	if err != nil {
		return
	}
	ses, err = app.storage.GetActiveSessionById(td.SessionId)
	if err != nil {
		return
	}
	err = verifyTokenData(td, ses)
	if err != nil {
		return
	}
	accessIsForbidden = false
	return
}

// Seeks the Marker Hash Sum inside the Token, searches the Storage for the
// Session having this Hash Sum, from that Session reads the Token Key which is
// returned as an Array (Slice) of Bytes. This Method is used together with the
// 'Parse' Method of a JWT-parsing Library and complies with the 'jwt.Keyfunc'
// Signature.
func (app *Application) getTokenKey(
	token *jwt.Token,
) (key interface{}, err error) {
	if token == nil {
		err = errors.New(ErrNullPointer)
		return
	}
	var claims jwt.MapClaims
	claims, err = jwtHelper.GetTokenClaims(token)
	if err != nil {
		return
	}
	var markerHash string
	markerHash, err = jwtHelper.GetTokenClaimMarkerHash(claims)
	if err != nil {
		return
	}
	key, err = app.storage.GetTokenKeyByMarkerHash(markerHash)
	if err != nil {
		return
	}
	var keyStr string
	var ok bool
	keyStr, ok = key.(string)
	if !ok {
		err = errors.New(ErrTypeCastFailure)
		return
	}
	key, err = hex.DecodeString(keyStr)
	if err != nil {
		return
	}

	return
}

// Verifies the auxiliary Data used with a JSON Web Token in this Project.
func verifyTokenData(
	td *jwtHelper.TokenData,
	ses *session.Session,
) (err error) {
	if (td == nil) || (ses == nil) {
		err = errors.New(ErrNullPointer)
		return
	}
	if (td.SessionId != ses.Id) ||
		(td.UniqueMarker != ses.Marker) ||
		(td.UniqueMarkerHash != ses.MarkerHash) {
		err = errors.New(ErrTokenDataMismatch)
		return
	}
	return
}
