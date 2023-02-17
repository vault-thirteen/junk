package httpserver

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
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	"github.com/vault-thirteen/junk/gfe/internal/kafka"
	"github.com/vault-thirteen/junk/gfe/internal/message"
	storageInterface "github.com/vault-thirteen/junk/gfe/pkg/repository"
)

// MsgFDebugConfig -- формат сообщения для отладки настроек.
const MsgFDebugConfig = "metrics http server configuration: %+v"

// HttpServer -- HTTP сервер для обработки запросов бизнес логики.
type HttpServer struct {
	// Внешние объекты.
	// Обнуление и изменение этих объектов запрещено.
	logger                 *zerolog.Logger
	kafka                  *kafka.Kafka
	storage                storageInterface.Storage
	applicationQuitSignals chan os.Signal
	metricsRegistry        *prometheus.Registry

	config *config.HttpServer
	server *http.Server
}

// NewHttpServer -- конструктор HTTP сервера, обрабатывающего запросы метрик.
func NewHttpServer(
	logger *zerolog.Logger,
	kafka *kafka.Kafka,
	storage storageInterface.Storage,
	applicationQuitSignals chan os.Signal,
	metricsRegistry *prometheus.Registry,
) (*HttpServer, error) {
	s := new(HttpServer)

	// Сохранение указателей на внешние объекты.
	s.logger = logger
	s.kafka = kafka
	s.storage = storage
	s.applicationQuitSignals = applicationQuitSignals
	s.metricsRegistry = metricsRegistry

	err := s.init()
	if err != nil {
		return nil, err
	}

	s.logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, s.config))

	return s, nil
}

// init производит первичную настройку HTTP сервера.
func (hs *HttpServer) init() (err error) {
	hs.config, err = config.GetSystemHttpServerConfig()
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

// initHttpRouter настраивает маршруты HTTP сервера.
func (hs *HttpServer) initHttpRouter() (httpRouter http.Handler, err error) {
	var router = httprouter.New()

	// 1. Обработчики доступности сервиса.

	// Примечание.
	// Liveness-хендлер расположен в HTTP сервере бизнес логике.

	// 1.1. Статус готовности сервиса (Readiness).
	// Данный путь зависит от настройки 'readinessProbe' в Kubernetes.
	router.GET("/ready", hs.handlerReadiness)

	// 2. Метрики.
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

// Start запускает HTTP сервер.
func (hs *HttpServer) Start() (err error) {
	msg := message.ComposeMessageWithPrefix(message.MsgPrefixMetrics, message.MsgHttpServerStarting)
	hs.logger.Info().Msg(msg)

	go func() {
		err = hs.server.ListenAndServe()
		if err != nil {
			msg = message.ComposeMessageWithPrefix(message.MsgPrefixMetrics, message.MsgHttpServerError)
			hs.logger.Err(err).Msg(msg)

			hs.applicationQuitSignals <- os.Interrupt
		}

		msg = message.ComposeMessageWithPrefix(message.MsgPrefixMetrics, message.MsgHttpServerStopped)
		hs.logger.Info().Msg(msg)
	}()

	return nil
}

// Stop останавливает HTTP сервер.
func (hs *HttpServer) Stop() (err error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*config.HttpServerShutdownTimeoutSec)
	defer cancelFn()

	err = hs.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
