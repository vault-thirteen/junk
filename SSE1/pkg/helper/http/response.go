package http

import (
	"net/http"

	"github.com/vault-thirteen/MIME"
	"github.com/vault-thirteen/header"
	"github.com/vault-thirteen/junk/SSE1/pkg/interfaces/response"
)

// Writes an Error to the HTTP Response Stream.
func SendHttpResponseError(
	w http.ResponseWriter,
	httpStatusCode int,
	errorText error,
) (err error) {
	w.WriteHeader(httpStatusCode)
	_, err = w.Write([]byte(errorText.Error()))
	if err != nil {
		return
	}
	return
}

// Writes an 'Internal Server' Error to the HTTP Response Stream.
func SendHttpResponseInternalServerError(
	w http.ResponseWriter,
	errorText error,
) (err error) {
	return SendHttpResponseError(w, http.StatusInternalServerError, errorText)
}

// Writes a 'Bad Request' Error to the HTTP Response Stream.
func SendHttpResponseBadRequestError(
	w http.ResponseWriter,
	errorText error,
) (err error) {
	return SendHttpResponseError(w, http.StatusBadRequest, errorText)
}

// Writes a 'Forbidden' Error to the HTTP Response Stream.
func SendHttpResponseForbiddenError(
	w http.ResponseWriter,
	errorText error,
) (err error) {
	return SendHttpResponseError(w, http.StatusForbidden, errorText)
}

// Sets an HTTP Header of the HTTP Response.
func SetHttpResponseHeader(
	w http.ResponseWriter,
	headerName string,
	headerValue string,
) {
	w.Header().Set(headerName, headerValue)
}

// Sets a secure HTTP Cookie of the HTTP Response.
func SetHttpSecureCookie(
	w http.ResponseWriter,
	cookieName string,
	cookieValue string,
	cookiePath string,
) {
	var cookie = http.Cookie{
		Name:     cookieName,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     cookiePath,
		Value:    cookieValue,
	}
	w.Header().Add(header.HttpHeaderSetCookie, cookie.String())
}

// Writes an Object in JSON Format to the HTTP Response Stream.
func RespondWithJsonObject(
	w http.ResponseWriter,
	responseObject response.IResponseObject,
) (err error) {
	var buffer []byte
	buffer, err = responseObject.MarshalJSON()
	if err != nil {
		return
	}
	w.Header().Set(
		header.HttpHeaderContentType,
		mime.TypeApplicationJson,
	)
	_, err = w.Write(buffer)
	if err != nil {
		return
	}
	return
}
