package c

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	hm "github.com/vault-thirteen/TR1/src/models/http"
	rm "github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

// Static files.
const (
	StaticFile_ApiJs      = "api.js"
	StaticFile_Argon2Js   = "argon2.js"
	StaticFile_Argon2Wasm = "argon2.wasm"
	StaticFile_BppJs      = "bpp.js"
	StaticFile_FaviconPng = "favicon.png"
	StaticFile_IndexHtml  = "index.html"
	StaticFile_LoaderJs   = "loader.js"
	StaticFile_StylesCss  = "styles.css"
)

// URL paths.
const (
	UrlPath_Root     = `/`
	UrlPath_Api      = `/api`
	UrlPath_Captcha  = `/captcha`
	UrlPath_Settings = `/settings`

	UrlPath_ApiJs      = `/api.js`
	UrlPath_Argon2Js   = `/argon2.js`
	UrlPath_Argon2Wasm = `/argon2.wasm`
	UrlPath_BppJs      = `/bpp.js`
	UrlPath_FaviconPng = `/favicon.png`
	UrlPath_IndexHtml  = `/index.html`
	UrlPath_LoaderJs   = `/loader.js`
	UrlPath_StylesCss  = `/styles.css`
)

// Errors.
const (
	ErrFUnknownRpcErrorCode = "unknown RPC error code: %v"
)

func (c *Controller) initGatewayRouter() {
	c.far.httpServer.SetHttpRouter(http.Handler(http.HandlerFunc(c.gatewayRouter)))
}

func (c *Controller) processInternalServerError(rw http.ResponseWriter, err error) {
	c.logError(err)
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusInternalServerError)
}
func (c *Controller) processRpcError(re *jrm1.RpcError, rw http.ResponseWriter) {
	httpStatusCode, ok := c.httpStatusCodesByRpcErrorCode[re.Code.Int()]
	if !ok {
		err := fmt.Errorf(ErrFUnknownRpcErrorCode, re.Code.Int())
		c.processInternalServerError(rw, err)
		return
	}

	switch httpStatusCode {
	case http.StatusInternalServerError:
		err := re.AsError()
		c.processInternalServerError(rw, err)
		return
	}

	c.respondWithPlainText(rw, re.AsError().Error(), httpStatusCode)
	return
}
func (c *Controller) respondWithPlainText(rw http.ResponseWriter, text string, httpStatusCode int) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.Header().Set(header.HttpHeaderContentType, hm.ContentType_PlainText)
	rw.WriteHeader(httpStatusCode)

	_, err := rw.Write([]byte(text))
	if err != nil {
		c.logError(err)
		return
	}
}
func (c *Controller) respondWithJsonObject(rw http.ResponseWriter, obj any) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.Header().Set(header.HttpHeaderContentType, hm.ContentType_Json)

	err := json.NewEncoder(rw).Encode(obj)
	if err != nil {
		c.logError(err)
		return
	}
}
func (c *Controller) respondBadRequest(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusBadRequest)
}
func (c *Controller) respondForbidden(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusForbidden)
}
func (c *Controller) respondMethodNotAllowed(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}
func (c *Controller) respondNotAcceptable(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusNotAcceptable)
}
func (c *Controller) respondNotFound(rw http.ResponseWriter) {
	if c.far.isDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, c.far.devModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.WriteHeader(http.StatusNotFound)
}

func (c *Controller) setTokenCookie(rw http.ResponseWriter, token string) {
	var cookie = &http.Cookie{
		Name:   rm.CookieName_Token,
		Value:  token,
		MaxAge: c.far.sessionMaxDuration,

		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}

	hm.SetCookie(rw, cookie)
}
func (c *Controller) clearTokenCookie(rw http.ResponseWriter) {
	var cookie = &http.Cookie{
		Name: rm.CookieName_Token,
		//Value
		//MaxAge

		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}

	hm.UnsetCookie(rw, cookie)
}

func (c *Controller) getClientIPAddress(req *http.Request) (cipa string, err error) {
	var host string

	if len(c.far.clientIPAddressSource_CustomHeader) == 0 {
		host, _, err = net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return "", err
		}

		return host, nil
	}

	host, err = hh.GetSingleHttpHeader(req, c.far.clientIPAddressSource_CustomHeader)
	if err != nil {
		return "", err
	}

	return host, nil
}

func (c *Controller) gatewayRouter(rw http.ResponseWriter, req *http.Request) {
	clientIPA, err := c.getClientIPAddress(req)
	if err != nil {
		c.processInternalServerError(rw, err)
		return
	}

	switch req.URL.Path {
	case UrlPath_Api:
		c.handleApiRequest(rw, req, clientIPA)
		return

	case UrlPath_Settings:
		c.handleSettingsRequest(rw, req)
		return

	case UrlPath_Root,
		UrlPath_IndexHtml:
		c.handleStaticFile(rw, StaticFile_IndexHtml, hm.ContentType_HtmlPage)
		return

	case UrlPath_StylesCss:
		c.handleStaticFile(rw, StaticFile_StylesCss, hm.ContentType_CssStyle)
		return

	case UrlPath_FaviconPng:
		c.handleStaticFile(rw, StaticFile_FaviconPng, hm.ContentType_PNG)
		return

	// JavaScript scripts.

	case UrlPath_ApiJs:
		c.handleStaticFile(rw, StaticFile_ApiJs, hm.ContentType_JavaScript)
		return

	case UrlPath_Argon2Js:
		c.handleStaticFile(rw, StaticFile_Argon2Js, hm.ContentType_JavaScript)
		return

	case UrlPath_BppJs:
		c.handleStaticFile(rw, StaticFile_BppJs, hm.ContentType_JavaScript)
		return

	case UrlPath_LoaderJs:
		c.handleStaticFile(rw, StaticFile_LoaderJs, hm.ContentType_JavaScript)
		return

	case UrlPath_Argon2Wasm:
		c.handleStaticFile(rw, StaticFile_Argon2Wasm, hm.ContentType_Wasm)
		return

	// Captcha.
	case UrlPath_Captcha:
		c.handleCaptchaRequest(rw, req)
		return

	default:
		c.handleStaticFile(rw, StaticFile_IndexHtml, hm.ContentType_HtmlPage)
		return
	}
}

func (c *Controller) handleCaptchaRequest(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		c.respondMethodNotAllowed(rw)
		return
	}

	c.far.captchaServiceProxy.Use(rw, req)

	return
}
func (c *Controller) handleStaticFile(rw http.ResponseWriter, fileName string, contentType string) {
	fileContents, err := c.far.fileServer.GetFile(fileName)
	if err != nil {
		c.logError(err)
		return
	}

	rw.Header().Set(header.HttpHeaderCacheControl, "max-age="+strconv.Itoa(c.far.cacheControlMaxAge))
	rw.Header().Set(header.HttpHeaderContentType, contentType)

	_, err = rw.Write(fileContents)
	if err != nil {
		c.logError(err)
		return
	}
}
