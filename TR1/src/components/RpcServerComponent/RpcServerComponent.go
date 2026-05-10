package rsc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

type RpcServerComponent struct {
	cfg        interfaces.IConfiguration
	errorsChan *chan error
	js         *jrm1.Processor
	listenDsn  string
	httpServer *http.Server
}

func (c *RpcServerComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	rsc := &RpcServerComponent{
		cfg:        cfg,
		errorsChan: controller.GetErrorsChan(),
	}

	// RPC processor.
	fnDur := rm.RpcDurationFieldName
	fnReqId := rm.RpcRequestIdFieldName

	ps := &jrm1.ProcessorSettings{
		CatchExceptions:    true,
		LogExceptions:      true,
		CountRequests:      true,
		DurationFieldName:  &fnDur,
		RequestIdFieldName: &fnReqId,
	}

	var js *jrm1.Processor
	js, err = jrm1.NewProcessor(ps)
	if err != nil {
		return nil, err
	}

	fns := controller.GetRpcFunctions()

	for _, fn := range fns {
		err = js.AddFunc(fn)
		if err != nil {
			return nil, err
		}
	}

	rsc.js = js

	// HTTP server.
	httpServerSettings := cfg.GetServer(cm.ServerType_Internal, cm.Protocol_HTTPS)
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
func (c *RpcServerComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *RpcServerComponent) httpRouter(rw http.ResponseWriter, req *http.Request) {
	c.js.ServeHTTP(rw, req)
}

func (c *RpcServerComponent) Start(s interfaces.IService) (err error) {
	go func() {
		httpServerSettings := c.cfg.GetServer(cm.ServerType_Internal, cm.Protocol_HTTPS)

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
func (c *RpcServerComponent) Stop(s interfaces.IService) (err error) {
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

func (c *RpcServerComponent) ReportStart() {
	fmt.Println(fmt.Sprintf("RpcServerComponent has started on %s", c.listenDsn))
}
func (c *RpcServerComponent) ReportStop() {
	fmt.Println("RpcServerComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *RpcServerComponent) {
	return x.(*RpcServerComponent)
}

// Non-standard methods.
