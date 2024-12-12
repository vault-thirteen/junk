package server

import (
	"encoding/json"
	http2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/http"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"log"
	"net/http"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/auxie/header"
)

// processInternalServerError logs the internal error and responds via HTTP.
func (srv *Server) processInternalServerError(rw http.ResponseWriter, err error) {
	srv.logError(err)
	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.WriteHeader(http.StatusInternalServerError)
}

// processRpcError finds an appropriate HTTP status code and message text for
// an RPC error and responds via HTTP.
func (srv *Server) processRpcError(moduleId int, re *jrm1.RpcError, rw http.ResponseWriter) {
	var httpStatusCode int
	var err error
	httpStatusCode, err = srv.getHttpStatusCodeByRpcErrorCode(moduleId, re.Code.Int())
	if err != nil {
		srv.processInternalServerError(rw, err)
		return
	}

	switch httpStatusCode {
	case http.StatusInternalServerError:
		srv.processInternalServerError(rw, err)
		return
	}

	srv.respondWithPlainText(rw, re.AsError().Error(), httpStatusCode)
	return
}

// respondWithPlainText responds via HTTP with a simple text message.
func (srv *Server) respondWithPlainText(rw http.ResponseWriter, text string, httpStatusCode int) {
	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.Header().Set(header.HttpHeaderContentType, http2.ContentType_PlainText)
	rw.WriteHeader(httpStatusCode)

	_, err := rw.Write([]byte(text))
	if err != nil {
		log.Println(err.Error())
		return
	}
}

// respondWithJsonObject responds via HTTP with a JSON object.
func (srv *Server) respondWithJsonObject(rw http.ResponseWriter, obj any) {
	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.Header().Set(header.HttpHeaderContentType, http2.ContentType_Json)

	err := json.NewEncoder(rw).Encode(obj)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (srv *Server) respondBadRequest(rw http.ResponseWriter) {
	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.WriteHeader(http.StatusBadRequest)
}

func (srv *Server) respondForbidden(rw http.ResponseWriter) {
	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.WriteHeader(http.StatusForbidden)
}

func (srv *Server) respondNotFound(rw http.ResponseWriter) {
	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.WriteHeader(http.StatusNotFound)
}

func (srv *Server) respondMethodNotAllowed(rw http.ResponseWriter) {
	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (srv *Server) respondNotAcceptable(rw http.ResponseWriter) {
	if srv.settings.GetSystemSettings().GetIsDeveloperMode() {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.GetSystemSettings().GetDevModeHttpHeaderAccessControlAllowOrigin())
	}
	rw.WriteHeader(http.StatusNotAcceptable)
}

func (srv *Server) setTokenCookie(rw http.ResponseWriter, token simple.WebTokenString) {
	var c = &http.Cookie{
		Name:   simple.CookieName_Token,
		Value:  token.ToString(),
		MaxAge: int(srv.settings.GetSystemSettings().GetSessionMaxDuration()),

		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}

	http2.SetCookie(rw, c)
}

func (srv *Server) clearTokenCookie(rw http.ResponseWriter) {
	var c = &http.Cookie{
		Name: simple.CookieName_Token,
		//Value
		//MaxAge

		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}

	http2.UnsetCookie(rw, c)
}
