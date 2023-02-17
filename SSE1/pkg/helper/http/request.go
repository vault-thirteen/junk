package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/vault-thirteen/errorz"
)

// Errors.
const (
	ErrfHeaderNotFound = "HTTP Header is not found: %v"
)

// Reads an HTTP Body of the HTTP Request.
func GetHttpRequestBody(
	r *http.Request,
) (reuestBody []byte, err error) {
	reuestBody, err = io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer func() {
		var derr = r.Body.Close()
		if derr != nil {
			err = errorz.Combine(err, derr)
		}
	}()
	return
}

// Reads an HTTP Cookie of the HTTP Request.
func GetHttpCookie(
	r *http.Request,
	cookieName string,
) (cookieValue string, err error) {
	var cookie *http.Cookie
	cookie, err = r.Cookie(cookieName)
	if err != nil {
		return
	}
	cookieValue = cookie.Value
	return
}

// Reads an HTTP Header of the HTTP Request.
func GetHttpRequestHeader(
	r *http.Request,
	headerName string,
) (headerValue string, err error) {
	var ok bool
	_, ok = r.Header[headerName]
	if !ok {
		err = fmt.Errorf(ErrfHeaderNotFound, headerName)
		return
	}
	headerValue = r.Header.Get(headerName)
	return
}
