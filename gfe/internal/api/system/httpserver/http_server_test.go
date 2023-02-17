package httpserver

import (
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	"github.com/vault-thirteen/junk/gfe/internal/kafka"
	iPrometheus "github.com/vault-thirteen/junk/gfe/internal/prometheus"
	iStorage "github.com/vault-thirteen/junk/gfe/internal/storage"
	storageInterface "github.com/vault-thirteen/junk/gfe/pkg/repository"
)

const (
	EnvVarNameHost = "GFE_METRICS_HTTP_SERVER_HOST"
	EnvVarNamePort = "GFE_METRICS_HTTP_SERVER_PORT"
)

func TestNewHttpServer(t *testing.T) {
	// Внимание!
	// Нельзя задавать переменным окружения названия тех переменных, которые
	// уже используются в операционной системе.

	// Arrange.
	var (
		logger                 *zerolog.Logger          = new(zerolog.Logger)
		kafkaObject            *kafka.Kafka             = new(kafka.Kafka)
		storageObject          storageInterface.Storage = new(iStorage.Storage)
		applicationQuitSignals chan os.Signal           = make(chan os.Signal, 123)
		prmts                  *iPrometheus.Prometheus
		metricsRegistry        *prometheus.Registry
		err                    error
		httpServerActual       *HttpServer
		httpServerExpected     *HttpServer
	)

	prmts, err = iPrometheus.NewPrometheus()
	assert.NoError(t, err)

	metricsRegistry = prmts.GetRegistry()

	err = os.Setenv(EnvVarNameHost, "localhost")
	assert.NoError(t, err)
	err = os.Setenv(EnvVarNamePort, "9999")
	assert.NoError(t, err)

	httpServerExpected = &HttpServer{
		logger:                 logger,
		kafka:                  kafkaObject,
		storage:                storageObject,
		applicationQuitSignals: applicationQuitSignals,

		config: &config.HttpServer{
			HttpServerHost: "localhost",
			HttpServerPort: 9999,
		},
		server: nil, // Null, ибо мы не можем просто проверить это поле.
	}

	// Act.
	httpServerActual, err = NewHttpServer(
		logger,
		kafkaObject,
		storageObject,
		applicationQuitSignals,
		metricsRegistry,
	)

	// Assert.
	assert.NoError(t, err)

	// Проверяем поля, переданные в конструктор.
	assert.Equal(t, httpServerActual.logger, httpServerExpected.logger)
	assert.Equal(t, httpServerActual.kafka, httpServerExpected.kafka)
	assert.Equal(t, httpServerActual.storage, httpServerExpected.storage)
	assert.Equal(t, httpServerActual.applicationQuitSignals, httpServerExpected.applicationQuitSignals)

	// Проверяем поле 'config'.
	assert.NotEqual(t, nil, httpServerActual.config)
	assert.Equal(t, httpServerActual.config.HttpServerHost, httpServerExpected.config.HttpServerHost)
	assert.Equal(t, httpServerActual.config.HttpServerPort, httpServerExpected.config.HttpServerPort)

	// Проверяем поле 'server'.
	assert.NotEqual(t, nil, httpServerActual.server)
	assert.Equal(t, httpServerActual.server.Addr, "localhost:9999")

	// Очищаем О.С. от мусора после теста.
	err = os.Setenv(EnvVarNameHost, "")
	assert.NoError(t, err)
	err = os.Setenv(EnvVarNamePort, "")
	assert.NoError(t, err)
}

func TestHttpServer_Start(t *testing.T) {
	// Внимание!
	// Нельзя задавать переменным окружения названия тех переменных, которые
	// уже используются в операционной системе.

	// Arrange.
	var (
		logger                 *zerolog.Logger          = new(zerolog.Logger)
		kafkaObject            *kafka.Kafka             = new(kafka.Kafka)
		storageObject          storageInterface.Storage = new(iStorage.Storage)
		applicationQuitSignals chan os.Signal           = make(chan os.Signal, 123)
		prmts                  *iPrometheus.Prometheus
		metricsRegistry        *prometheus.Registry
		err                    error
		httpServerActual       *HttpServer
	)

	prmts, err = iPrometheus.NewPrometheus()
	assert.NoError(t, err)

	metricsRegistry = prmts.GetRegistry()

	err = os.Setenv(EnvVarNameHost, "localhost")
	assert.NoError(t, err)
	err = os.Setenv(EnvVarNamePort, "9999")
	assert.NoError(t, err)

	httpServerActual, err = NewHttpServer(
		logger,
		kafkaObject,
		storageObject,
		applicationQuitSignals,
		metricsRegistry,
	)
	assert.NoError(t, err)

	// Act.
	err = httpServerActual.Start()

	// Assert.
	assert.NoError(t, err)

	// Закрываем открытые ресурсы.
	err = httpServerActual.Stop()
	assert.NoError(t, err)

	// Очищаем О.С. от мусора после теста.
	err = os.Setenv(EnvVarNameHost, "")
	assert.NoError(t, err)
	err = os.Setenv(EnvVarNamePort, "")
	assert.NoError(t, err)
}

func TestHttpServer_Stop(t *testing.T) {
	// Внимание!
	// Нельзя задавать переменным окружения названия тех переменных, которые
	// уже используются в операционной системе.

	// Arrange.
	var (
		logger                 *zerolog.Logger          = new(zerolog.Logger)
		kafkaObject            *kafka.Kafka             = new(kafka.Kafka)
		storageObject          storageInterface.Storage = new(iStorage.Storage)
		applicationQuitSignals chan os.Signal           = make(chan os.Signal, 123)
		prmts                  *iPrometheus.Prometheus
		metricsRegistry        *prometheus.Registry
		err                    error
		httpServerActual       *HttpServer
	)

	prmts, err = iPrometheus.NewPrometheus()
	assert.NoError(t, err)

	metricsRegistry = prmts.GetRegistry()

	err = os.Setenv(EnvVarNameHost, "localhost")
	assert.NoError(t, err)
	err = os.Setenv(EnvVarNamePort, "9999")
	assert.NoError(t, err)

	httpServerActual, err = NewHttpServer(
		logger,
		kafkaObject,
		storageObject,
		applicationQuitSignals,
		metricsRegistry,
	)
	assert.NoError(t, err)
	err = httpServerActual.Start()
	assert.NoError(t, err)
	time.Sleep(time.Second)

	// Act.
	err = httpServerActual.Stop()

	// Assert.
	assert.NoError(t, err)

	// Очищаем О.С. от мусора после теста.
	err = os.Setenv(EnvVarNameHost, "")
	assert.NoError(t, err)
	err = os.Setenv(EnvVarNamePort, "")
	assert.NoError(t, err)
}
