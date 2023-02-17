package application

import (
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/kr/pretty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/vault-thirteen/junk/SSE2/internal/api/hs"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/fsl"
	inKafka "github.com/vault-thirteen/junk/SSE2/internal/application/component/kafka/input"
	outKafka "github.com/vault-thirteen/junk/SSE2/internal/application/component/kafka/output"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/metrics"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/service"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
)

const MsgFDebugConfig = "storage configuration: %+v"

func (a *Application) init() (err error) {
	err = a.initService()
	if err != nil {
		return err
	}

	err = a.initLogger()
	if err != nil {
		return err
	}

	err = a.initControls()
	if err != nil {
		return err
	}

	err = a.initPrometheus()
	if err != nil {
		return err
	}

	err = a.initStorageSettings()
	if err != nil {
		return err
	}

	err = a.initFileSizeLimiter()
	if err != nil {
		return err
	}

	err = a.initKafka()
	if err != nil {
		return err
	}

	err = a.initHttpServer()
	if err != nil {
		return err
	}

	err = a.service.Configure(
		a.logger,
		a.fileSizeLimiter,
		a.inputKafka,
		a.outputKafka,
		a.storageSettings,
		a.prometheus.GetMetrics(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) initService() (err error) {
	a.service, err = service.NewService()
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) initControls() (err error) {
	const QuitSignalsSize = 64

	a.quitSignals = make(chan os.Signal, QuitSignalsSize)
	signal.Notify(a.quitSignals, os.Interrupt)

	a.isStarted.Store(false)

	return nil
}

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

func (a *Application) initPrometheus() (err error) {
	a.prometheus, err = metrics.NewPrometheus()
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) initKafka() (err error) {
	a.inputKafka, err = inKafka.NewKafka(
		a.logger,
		a.prometheus.GetMetrics(),
		a.service,
	)
	if err != nil {
		return err
	}

	a.outputKafka, err = outKafka.NewKafka(
		a.logger,
		a.prometheus.GetMetrics(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) initStorageSettings() (err error) {
	a.storageSettings, err = config.GetStorageSettings()
	if err != nil {
		return err
	}

	a.logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, a.storageSettings))

	return nil
}

func (a *Application) initFileSizeLimiter() (err error) {
	a.fileSizeLimiter, err = fsl.NewFileSizeLimiter(
		a.logger,
		a.service.GetFileSizeLimitSettingsFile(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) initHttpServer() (err error) {
	a.httpServer, err = hs.NewHttpServer(
		a.logger,
		a.service,
		a.quitSignals,
		a.prometheus.GetRegistry(),
	)
	if err != nil {
		return err
	}

	return nil
}
