package hs

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kr/pretty"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
	"github.com/vault-thirteen/junk/SSE2/internal/messages"
	serviceInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/service"
)

const MsgFDebugConfig = "metrics http server configuration: %+v"

type HttpServer struct {
	logger                 *zerolog.Logger
	service                serviceInterface.Service
	applicationQuitSignals chan os.Signal
	metricsRegistry        *prometheus.Registry

	config *config.HttpServer
	server *http.Server
}

func NewHttpServer(
	logger *zerolog.Logger,
	service serviceInterface.Service,
	applicationQuitSignals chan os.Signal,
	metricsRegistry *prometheus.Registry,
) (*HttpServer, error) {
	s := new(HttpServer)

	s.logger = logger
	s.service = service
	s.applicationQuitSignals = applicationQuitSignals
	s.metricsRegistry = metricsRegistry

	err := s.init()
	if err != nil {
		return nil, err
	}

	s.logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, s.config))

	return s, nil
}

func (hs *HttpServer) init() (err error) {
	hs.config, err = config.GetHttpServerConfig()
	if err != nil {
		return err
	}

	hs.server = &http.Server{
		Addr: net.JoinHostPort(
			hs.config.HttpServerHost,
			strconv.FormatUint(uint64(hs.config.HttpServerPort), 10),
		),
	}

	hs.server.Handler, err = hs.initHttpRouter()
	if err != nil {
		return err
	}

	return nil
}

func (hs *HttpServer) initHttpRouter() (httpRouter http.Handler, err error) {
	var router = httprouter.New()

	router.GET("/live", hs.handlerLiveness)

	router.GET("/ready", hs.handlerReadiness)

	router.Handler(
		http.MethodGet,
		"/metrics",
		promhttp.InstrumentMetricHandler(
			hs.metricsRegistry, promhttp.HandlerFor(
				hs.metricsRegistry,
				promhttp.HandlerOpts{
					ErrorLog:      log.Default(),
					ErrorHandling: promhttp.HTTPErrorOnError,
					Registry:      hs.metricsRegistry,
					Timeout:       time.Second * 60,
				},
			),
		),
	)

	return router, nil
}

func (hs *HttpServer) Start() (err error) {
	hs.logger.Info().Msg(messages.MsgHttpServerStarting)

	go func() {
		err = hs.server.ListenAndServe()
		if err != nil {
			hs.logger.Err(err).Msg(messages.MsgHttpServerError)

			hs.applicationQuitSignals <- os.Interrupt
		}

		hs.logger.Info().Msg(messages.MsgHttpServerStopped)
	}()

	return nil
}

func (hs *HttpServer) Stop() (err error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*config.HttpServerShutdownTimeoutSec)
	defer cancelFn()

	err = hs.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
