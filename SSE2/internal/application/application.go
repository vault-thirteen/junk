package application

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/api/hs"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/fsl"
	"github.com/vault-thirteen/junk/SSE2/internal/application/component/metrics"
	"github.com/vault-thirteen/junk/SSE2/internal/application/config"
	"github.com/vault-thirteen/junk/SSE2/internal/messages"
	inputKafkaInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/kafka/input"
	outputKafkaInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/kafka/output"
	serviceInterface "github.com/vault-thirteen/junk/SSE2/pkg/interfaces/service"
	"go.uber.org/atomic"
)

type Application struct {
	loggerConfig *config.Logger
	logger       *zerolog.Logger

	quitSignals chan os.Signal
	isStarted   atomic.Bool
	starterLock sync.Mutex
	stopperLock sync.Mutex

	prometheus      *metrics.Prometheus
	inputKafka      inputKafkaInterface.Kafka
	outputKafka     outputKafkaInterface.Kafka
	httpServer      *hs.HttpServer
	storageSettings *config.Storage
	fileSizeLimiter *fsl.FileSizeLimiter
	service         serviceInterface.Service
}

func NewApplication() (app *Application, err error) {
	app = new(Application)

	err = app.init()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *Application) Start() (err error) {
	a.starterLock.Lock()
	defer a.starterLock.Unlock()
	if a.isStarted.Load() {
		return messages.ErrIsAlreadyStarted
	}

	a.inputKafka.Run()
	a.inputKafka.Wait()

	err = a.service.Start()
	if err != nil {
		return err
	}

	go a.readServiceErrors()

	err = a.httpServer.Start()
	if err != nil {
		return err
	}

	a.isStarted.Store(true)

	return nil
}

func (a *Application) WaitForQuitSignal() (err error) {
	sig := <-a.quitSignals
	a.logger.Info().Msgf(messages.MsgFSignalReceived, sig)

	err = a.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) MustBeNoError(err error) {
	if err != nil {
		a.logger.Err(err).Send()
		os.Exit(config.ExitCodeOnError)
	}
}

func (a *Application) Stop() (err error) {
	a.stopperLock.Lock()
	defer a.stopperLock.Unlock()
	if !a.isStarted.Load() {
		return messages.ErrIsNotStarted
	}

	err = a.inputKafka.Close()
	if err != nil {
		return err
	}

	err = a.service.Stop()
	if err != nil {
		return err
	}

	err = a.outputKafka.Close()
	if err != nil {
		return err
	}

	err = a.stopHttpServers()
	if err != nil {
		return err
	}

	a.isStarted.Store(false)

	return nil
}

func (a *Application) stopHttpServers() (err error) {
	const DelayAfterHttpServerShutdownMs = 100

	err = a.httpServer.Stop()
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * DelayAfterHttpServerShutdownMs)

	return nil
}

func (a *Application) readServiceErrors() {
	errorsChannel := a.service.GetErrorsChannel()

	a.logger.Info().Msg(messages.MsgServiceErrorReaderStart)

	for err := range errorsChannel {
		a.logger.Err(err).Send()
	}

	a.logger.Info().Msg(messages.MsgServiceErrorReaderStop)
}
