package rmp

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	ccp "github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	ccse "github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
)

type Proxy struct {
	shortName string
	proxy     *httputil.ReverseProxy
}

func NewProxyFromSettings(settings *ccse.CommonConfigurationServiceEntry, shortName string) (proxy *Proxy, err error) {
	schema := settings.GetParameterAsString(ccp.Schema)
	host := settings.GetParameterAsString(ccp.Host)
	port := settings.GetParameterAsInt(ccp.Port)
	path := settings.GetParameterAsString(ccp.Path)
	targetAddr := schema + "://" + net.JoinHostPort(host, strconv.Itoa(port)) + path

	var targetUrl *url.URL
	targetUrl, err = url.Parse(targetAddr)
	if err != nil {
		return nil, err
	}

	var rp *httputil.ReverseProxy
	rp = httputil.NewSingleHostReverseProxy(targetUrl)

	proxy = &Proxy{
		shortName: shortName,
		proxy:     rp,
	}

	return proxy, nil
}

func (p *Proxy) Use(rw http.ResponseWriter, req *http.Request) {
	p.proxy.ServeHTTP(rw, req)
}
