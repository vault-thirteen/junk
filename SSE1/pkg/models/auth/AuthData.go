package auth

import (
	jwtHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/jwt"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/http/request"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/session"
)

// Authentication Data.
type AuthData struct {
	Session   *session.Session
	TokenData *jwtHelper.TokenData
	Machine   *request.UserLogRequestMachine
}
