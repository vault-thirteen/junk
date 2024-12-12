package simple

import (
	"net/http"

	hh "github.com/vault-thirteen/auxie/http-helper"
)

const (
	CookieName_Token = "token"
)

// GetToken tries to read a token. If a token is not found, null is returned.
func GetToken(req *http.Request) (token *WebTokenString, err error) {
	var cookie *http.Cookie
	cookie, err = hh.GetCookieByName(req, CookieName_Token)
	if err != nil {
		return nil, err
	}

	if cookie == nil {
		return nil, nil
	}

	token = new(WebTokenString)
	*token = WebTokenString(cookie.Value)

	return token, nil
}
