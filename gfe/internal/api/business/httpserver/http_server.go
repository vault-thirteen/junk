package httpserver

import (
	"context"
	"crypto/rsa"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kr/pretty"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	"github.com/vault-thirteen/junk/gfe/internal/kafka"
	"github.com/vault-thirteen/junk/gfe/internal/message"
	"github.com/vault-thirteen/junk/gfe/internal/prometheus"
	"github.com/vault-thirteen/junk/gfe/pkg/repository"
)

// MsgFDebugConfig -- формат сообщения для отладки настроек.
const MsgFDebugConfig = "business logics http server configuration: %+v"

// HttpServer -- HTTP сервер для обработки запросов бизнес логики.
type HttpServer struct {
	// Внешние объекты.
	// Обнуление и изменение этих объектов запрещено.
	logger                 *zerolog.Logger
	kafka                  *kafka.Kafka
	storage                repository.Storage
	jwtRsaPublicKey        *rsa.PublicKey
	prometheusMetrics      *prometheus.Metrics
	applicationQuitSignals chan os.Signal

	config *config.HttpServer
	server *http.Server
}

// NewHttpServer -- конструктор HTTP сервера, обрабатывающего запросы бизнес
// логики.
func NewHttpServer(
	logger *zerolog.Logger,
	kafka *kafka.Kafka,
	storage repository.Storage,
	jwtRsaPublicKey *rsa.PublicKey,
	prometheusMetrics *prometheus.Metrics,
	applicationQuitSignals chan os.Signal,
) (*HttpServer, error) {
	s := new(HttpServer)

	// Сохранение указателей на внешние объекты.
	s.logger = logger
	s.kafka = kafka
	s.storage = storage
	s.jwtRsaPublicKey = jwtRsaPublicKey
	s.prometheusMetrics = prometheusMetrics
	s.applicationQuitSignals = applicationQuitSignals

	err := s.init()
	if err != nil {
		return nil, err
	}

	s.logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, s.config))

	return s, nil
}

// init производит первичную настройку HTTP сервера.
func (hs *HttpServer) init() (err error) {
	hs.config, err = config.GetBusinessHttpServerConfig()
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
	// Readiness-хендлер расположен в HTTP сервере метрик.

	// 1.1. Статус жизни сервиса (Liveness).
	// Данный путь зависит от настройки 'livenessProbe' в Kubernetes.
	router.GET("/live", hs.handlerLiveness)

	// 2. Обработчики запросов на получение истории событий.
	// Примечание. Для склеиваемых (агрегируемых) событий -- датой
	// виртуального склеенного события является самая поздняя (максимальная)
	// дата событий, входящих в склейку.

	// 2.1. Получение списка всех событий по файлу.
	router.GET(
		"/file-events/all",
		hs.middlewareAuthentication(hs.handlerGetAllEventsForFile),
	)

	// 2.2. Получение списка нескольких последних (недавних) событий по файлу.
	router.GET(
		"/file-events/last-n",
		hs.middlewareAuthentication(hs.handlerGetLastNEventsForFile),
	)

	// 2.3. Получение списка типов событий.
	router.GET(
		"/file-event/types",
		hs.middlewareAuthentication(hs.handlerGetFileEventTypes),
	)

	return router, nil
}

// Start запускает HTTP сервер.
func (hs *HttpServer) Start() (err error) {
	msg := message.ComposeMessageWithPrefix(message.MsgPrefixBusinessLogics, message.MsgHttpServerStarting)
	hs.logger.Info().Msg(msg)

	go func() {
		err = hs.server.ListenAndServe()
		if err != nil {
			msg = message.ComposeMessageWithPrefix(message.MsgPrefixBusinessLogics, message.MsgHttpServerError)
			hs.logger.Err(err).Msg(msg)

			hs.applicationQuitSignals <- os.Interrupt
		}

		msg = message.ComposeMessageWithPrefix(message.MsgPrefixBusinessLogics, message.MsgHttpServerStopped)
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
