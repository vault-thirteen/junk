package application

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
	bhs "github.com/vault-thirteen/junk/gfe/internal/api/business/httpserver"
	shs "github.com/vault-thirteen/junk/gfe/internal/api/system/httpserver"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	"github.com/vault-thirteen/junk/gfe/internal/jwt"
	"github.com/vault-thirteen/junk/gfe/internal/kafka"
	"github.com/vault-thirteen/junk/gfe/internal/message"
	"github.com/vault-thirteen/junk/gfe/internal/prometheus"
	storageInterface "github.com/vault-thirteen/junk/gfe/pkg/repository"
	"go.uber.org/atomic"
)

// Application -- приложение.
type Application struct {
	// 1. Ведущий журнала.
	loggerConfig *config.Logger
	logger       *zerolog.Logger

	// 2. Структуры управления приложением:
	//	- Канал для входящих сигналов от О.С. о завершении;
	//	- Флаг, показывающий что приложение запущено;
	//	- Мьютексы для защиты от неправильного запуска и остановки.
	quitSignals chan os.Signal
	isStarted   atomic.Bool
	starterLock sync.Mutex
	stopperLock sync.Mutex

	// 3. Метрики Prometheus.
	prometheus *prometheus.Prometheus

	// 4. JWT инфраструктура:
	jwt *jwt.Jwt

	// 5. Kafka.
	kafka *kafka.Kafka

	// 6. Хранилище (база данных) и его коннектор.
	storage storageInterface.Storage

	// 7. HTTP сервер для бизнес логики.
	businessHttpServer *bhs.HttpServer

	// 8. HTTP сервер для метрик.
	metricsHttpServer *shs.HttpServer
}

// NewApplication -- создатель приложения.
// Основные действия происходят в методе 'init':
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
// Источник комментария -- метод 'init'.
func NewApplication() (app *Application, err error) {
	app = new(Application)

	err = app.init()
	if err != nil {
		return nil, err
	}

	return app, nil
}

// Start запускает приложение.
// Порядок запуска компонентов важен.
func (a *Application) Start() (err error) {
	// Защита от дурака.
	a.starterLock.Lock()
	defer a.starterLock.Unlock()
	if a.isStarted.Load() {
		return message.ErrIsAlreadyStarted
	}

	// Хранилище.
	a.storage.Open()
	a.storage.Wait()

	// Kafka.
	a.kafka.Run()
	a.kafka.Wait()

	// HTTP сервер для метрик.
	err = a.metricsHttpServer.Start()
	if err != nil {
		return err
	}

	// HTTP сервер для бизнес логики.
	// Поскольку на этом сервере висит хендлер Liveness, то запускаем его в
	// самую последнюю очередь.
	err = a.businessHttpServer.Start()
	if err != nil {
		return err
	}

	// Изменяем состояние.
	a.isStarted.Store(true)

	return nil
}

// WaitForQuitSignal ждёт сигнала завершения и запускает остановку приложения.
func (a *Application) WaitForQuitSignal() (err error) {
	sig := <-a.quitSignals
	a.logger.Info().Msgf(message.MsgFSignalReceived, sig)

	err = a.Stop()
	if err != nil {
		return err
	}

	return nil
}

// MustBeNoError выводит фатальную ошибку и завершает процесс ОС, если она есть.
func (a *Application) MustBeNoError(err error) {
	if err != nil {
		a.logger.Err(err).Send()
		os.Exit(config.ExitCodeOnError)
	}
}
