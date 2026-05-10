package hsc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

type HttpServerComponent struct {
	cfg        interfaces.IConfiguration
	errorsChan *chan error
	listenDsn  string
	httpServer *http.Server
}

func (c *HttpServerComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	rsc := &HttpServerComponent{
		cfg:        cfg,
		errorsChan: controller.GetErrorsChan(),
	}

	// HTTP server.
	httpServerSettings := cfg.GetServer(cm.ServerType_External, cm.Protocol_HTTPS)
	host := httpServerSettings.GetParameterAsString(ccp.Host)
	port := httpServerSettings.GetParameterAsInt(ccp.Port)
	listenDsn := net.JoinHostPort(host, strconv.Itoa(port))

	httpServer := &http.Server{
		Addr:    listenDsn,
		Handler: http.Handler(http.HandlerFunc(rsc.httpRouter)),
	}

	rsc.listenDsn = listenDsn
	rsc.httpServer = httpServer

	return rsc, nil
}
func (c *HttpServerComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *HttpServerComponent) httpRouter(rw http.ResponseWriter, req *http.Request) {
	// This is a default HTTP router. It provides very basic functionality. If
	// advanced functionality is needed, the default router must be changed
	// with the 'SetHttpRouter' method externally, i.e. from the controller.
	_, err := rw.Write([]byte("default router"))
	if err != nil {
		log.Println(err)
	}
}

func (c *HttpServerComponent) Start(s interfaces.IService) (err error) {
	go func() {
		httpServerSettings := c.cfg.GetServer(cm.ServerType_External, cm.Protocol_HTTPS)

		certFile := httpServerSettings.GetParameterAsString(ccp.CertFile)
		keyFile := httpServerSettings.GetParameterAsString(ccp.KeyFile)

		var listenError error
		listenError = c.httpServer.ListenAndServeTLS(certFile, keyFile)

		if (listenError != nil) && (!errors.Is(listenError, http.ErrServerClosed)) {
			*c.errorsChan <- listenError
		}
	}()

	return nil
}
func (c *HttpServerComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	err = c.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	c.ReportStop()

	return nil
}

func (c *HttpServerComponent) ReportStart() {
	fmt.Println(fmt.Sprintf("HttpServerComponent has started on %s", c.listenDsn))
}
func (c *HttpServerComponent) ReportStop() {
	fmt.Println("HttpServerComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *HttpServerComponent) {
	return x.(*HttpServerComponent)
}

// Non-standard methods.

func (c *HttpServerComponent) SetHttpRouter(router http.Handler) {
	c.httpServer.Handler = router
}
