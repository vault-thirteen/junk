package application

import (
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	bhs "github.com/vault-thirteen/junk/gfe/internal/api/business/httpserver"
	shs "github.com/vault-thirteen/junk/gfe/internal/api/system/httpserver"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	"github.com/vault-thirteen/junk/gfe/internal/jwt"
	"github.com/vault-thirteen/junk/gfe/internal/kafka"
	"github.com/vault-thirteen/junk/gfe/internal/prometheus"
	"github.com/vault-thirteen/junk/gfe/internal/storage"
)

// init производит первичную настройку приложения:
//   - настройка ведения журнала;
//   - настройка структур управления приложением;
//   - настройка метрик;
//   - настройка ключей для JWT веб токенов;
//   - настройка хранилища;
//   - настройка Kafka;
//   - настройка HTTP сервера для бизнес логики;
//   - настройка HTTP сервера для метрик.
//
// Примечание.
// Этот комментарий может дублироваться в других методах.
//
// Порядок инициализации компонентов важен.
func (a *Application) init() (err error) {
	// Инициализация :: Журнал.
	err = a.initLogger()
	if err != nil {
		return err
	}

	// Инициализация :: Структуры управления приложением.
	err = a.initControls()
	if err != nil {
		return err
	}

	// Инициализация :: Метрики.
	err = a.initPrometheus()
	if err != nil {
		return err
	}

	// Инициализация :: JWT инфраструктура.
	err = a.initJwt()
	if err != nil {
		return err
	}

	// Инициализация :: Хранилище (база данных).
	err = a.initStorage()
	if err != nil {
		return err
	}

	// Инициализация :: Kafka.
	err = a.initKafka()
	if err != nil {
		return err
	}

	// Инициализация :: HTTP сервер для бизнес логики.
	err = a.initBusinessHttpServer()
	if err != nil {
		return err
	}

	// Инициализация :: HTTP сервер для метрик.
	err = a.initMetricsHttpServer()
	if err != nil {
		return err
	}

	return nil
}

// initLogger производит первичную настройку ведения журнала.
func (a *Application) initLogger() (err error) {
	var loggerConfig *config.Logger
	loggerConfig, err = config.GetLoggerConfig()
	if err != nil {
		return err
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if loggerConfig.IsDebugEnabled {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		sarama.Logger = log.New(os.Stderr, "[Sarama] ", log.LstdFlags)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	logger := zerolog.New(os.Stderr).With().
		Timestamp().
		Stack().
		Logger()

	a.loggerConfig = loggerConfig
	a.logger = &logger

	return nil
}

// initControls производит первичную настройку управляющих структур.
func (a *Application) initControls() (err error) {
	const quitSignalsSize = 64

	a.quitSignals = make(chan os.Signal, quitSignalsSize)
	signal.Notify(a.quitSignals, os.Interrupt)

	a.isStarted.Store(false)

	return nil
}

// initPrometheus производит первичную настройку метрик Prometheus.
func (a *Application) initPrometheus() (err error) {
	a.prometheus, err = prometheus.NewPrometheus()
	if err != nil {
		return err
	}

	return nil
}

// initJwt производит первичную настройку JWT инфраструктуры.
func (a *Application) initJwt() (err error) {
	a.jwt, err = jwt.NewJwt(a.logger)
	if err != nil {
		return err
	}

	return nil
}

// initKafka производит первичную настройку Kafka.
func (a *Application) initKafka() (err error) {
	a.kafka, err = kafka.NewKafka(
		a.logger,
		a.storage,
		a.prometheus.GetMetrics(),
	)
	if err != nil {
		return err
	}

	return nil
}

// initStorage производит первичную настройку хранилища.
func (a *Application) initStorage() (err error) {
	a.storage, err = storage.NewStorage(
		a.logger,
		a.prometheus.GetMetrics(),
	)
	if err != nil {
		return err
	}

	return nil
}

// initBusinessHttpServer производит первичную настройку HTTP сервера,
// обрабатывающего запросы бизнес логики.
func (a *Application) initBusinessHttpServer() (err error) {
	a.businessHttpServer, err = bhs.NewHttpServer(
		a.logger,
		a.kafka,
		a.storage,
		a.jwt.GetRsaPublicKey(),
		a.prometheus.GetMetrics(),
		a.quitSignals,
	)
	if err != nil {
		return err
	}

	return nil
}

// initMetricsHttpServer производит первичную настройку HTTP сервера,
// обрабатывающего запросы метрик.
func (a *Application) initMetricsHttpServer() (err error) {
	a.metricsHttpServer, err = shs.NewHttpServer(
		a.logger,
		a.kafka,
		a.storage,
		a.quitSignals,
		a.prometheus.GetRegistry(),
	)
	if err != nil {
		return err
	}

	return nil
}
