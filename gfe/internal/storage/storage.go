package storage

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	iPrometheus "github.com/vault-thirteen/junk/gfe/internal/prometheus"
	storageInterface "github.com/vault-thirteen/junk/gfe/pkg/repository"
	"go.uber.org/atomic"
)

// Сообщения.
const (
	// MsgStorageNotReady -- сообщение "хранилище не готово".
	MsgStorageNotReady = "storage is not ready"

	msgStorageConnectorStart          = "storage connector has started"
	msgStorageNotConnected            = "storage is not connected"
	msgStorageConnected               = "storage is connected"
	msgStoragePingFailure             = "Ping failure"
	msgStorageConnecting              = "connecting to storage ..."
	msgStorageConnectionFailure       = "connection failure"
	msgStorageConnectorStop           = "storage connector has stopped"
	msgStoragePreparedStatementsReady = "storage prepared statements are ready"
)

// Ошибки.
const (
	// ErrStorageConnectionIsNull -- сообщение "соединение -- нуль".
	ErrStorageConnectionIsNull = "storage connection is null"
)

// ComponentsCount -- количество внутренних компонентов хранилища.
// В данный момент хранилище содержит следующие компоненты:
//  1. Соединятель (коннектор) базы данных.
const ComponentsCount = 1

// Индексы SQL запросов в списке запросов.
const (
	// SqlQueryIndexFirst -- индекс первой записи в массиве.
	// Всегда должен быть равен нулю.
	// Носит информационный характер, поэтому использования в коде нет.
	SqlQueryIndexFirst = 0

	// SqlQueryIndexLast -- индекс последней записи в массиве.
	// Зависит от количества записей.
	SqlQueryIndexLast = 4

	// SqlQueryIndexForQueryToGetFileEventTypes -- индекс SQL запроса на
	// выборку типов событий.
	SqlQueryIndexForQueryToGetFileEventTypes = 0

	// SqlQueryIndexForQueryToInsertFileSimpleEvent  -- индекс SQL запроса на
	// вставку события простого типа.
	SqlQueryIndexForQueryToInsertFileSimpleEvent = 1

	// SqlQueryIndexForQueryToInsertFileDownloadEvent -- индекс SQL запроса на
	// вставку события типа "Скачивание".
	SqlQueryIndexForQueryToInsertFileDownloadEvent = 2

	// SqlQueryIndexForQueryToSelectFileEventsAll -- индекс SQL запроса на
	// выборку всех событий по файлу.
	SqlQueryIndexForQueryToSelectFileEventsAll = 3

	// SqlQueryIndexForQueryToSelectFileEventsLimited -- индекс SQL запроса на
	// выборку событий по файлу с ограничителем количества.
	SqlQueryIndexForQueryToSelectFileEventsLimited = 4
)

// MsgFDebugConfig -- формат сообщения для отладки настроек.
const MsgFDebugConfig = "storage configuration: %+v"

// Storage -- хранилище (база дынных).
type Storage struct {
	// Внешние объекты.
	// Обнуление и изменение этих объектов запрещено.
	logger *zerolog.Logger

	// Канал для оповещения внутренних компонентов о надобности завершения.
	// Компоненты читают из этого канала сигналы о завершении.
	// Следовательно, для завершения всех компонентов, нужно отправлять N
	// сигналов, где N равно количеству компонентов.
	close chan bool

	// Количество внутренних компонентов хранилища.
	// Этот параметр важен для правильной остановки внутренних компонентов.
	componentsCount byte

	// Структура управления внутренним компонентом "соединятель".
	connectorWG *sync.WaitGroup

	// Флаг состояния готовности хранилища.
	isReady *atomic.Bool

	// Настройки.
	config *config.Storage

	// Соединение с базой данных.
	connection *sqlx.DB

	// Метрики.
	prometheusMetrics *iPrometheus.Metrics

	// SQL запросы.
	sqlQueryTexts                  []string
	sqlQueryPreparedStatements     []*sqlx.Stmt
	sqlQueryPreparedStatementsLock sync.RWMutex
}

// NewStorage -- конструктор хранилища.
// Автоматически запускает хранилище при создании.
func NewStorage(
	logger *zerolog.Logger,
	prometheusMetrics *iPrometheus.Metrics,
) (storageInterface.Storage, error) {
	s := new(Storage)

	// Сохранение указателей на внешние объекты.
	s.logger = logger
	s.prometheusMetrics = prometheusMetrics

	err := s.init()
	if err != nil {
		return nil, err
	}

	s.logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, s.config))

	return s, nil
}

// init производит первичную настройку.
func (s *Storage) init() (err error) {
	s.config, err = config.GetStorageConfig()
	if err != nil {
		return err
	}

	s.componentsCount = ComponentsCount
	s.close = make(chan bool, s.componentsCount)

	s.initSqlQueryTexts()

	s.connectorWG = new(sync.WaitGroup)

	s.isReady = new(atomic.Bool)
	s.isReady.Store(false)

	return nil
}

// prepareSqlQueryStatements подготавливает statements для SQL запросов.
func (s *Storage) prepareSqlQueryStatements() (err error) {
	if s.isConnectionNull() {
		return errors.New(msgStorageNotConnected)
	}

	s.sqlQueryPreparedStatementsLock.Lock()
	defer s.sqlQueryPreparedStatementsLock.Unlock()

	s.sqlQueryPreparedStatements = make([]*sqlx.Stmt, len(s.sqlQueryTexts))

	for i := range s.sqlQueryTexts {
		s.sqlQueryPreparedStatements[i], err = s.connection.Preparex(s.sqlQueryTexts[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// initSqlQueryTexts инициализирует тексты SQL запросов.
func (s *Storage) initSqlQueryTexts() {
	const sqlQueryTerminator = `;`

	s.sqlQueryTexts = make([]string, SqlQueryIndexLast+1)

	s.sqlQueryTexts[SqlQueryIndexForQueryToGetFileEventTypes] =
		`SELECT id, description_ru, description_en, is_simple, is_aggregated FROM public.event_types` +
			sqlQueryTerminator

	s.sqlQueryTexts[SqlQueryIndexForQueryToInsertFileSimpleEvent] =
		`INSERT INTO public.simple_events (user_id, file_id, event_type_id, event_time) VALUES ($1, $2, $3, $4)` +
			sqlQueryTerminator

	s.sqlQueryTexts[SqlQueryIndexForQueryToInsertFileDownloadEvent] =
		`INSERT INTO public.download_events (user_id, file_id, event_time) VALUES ($1, $2, $3)` +
			sqlQueryTerminator

	// Основа запросов для выборки событий по файлу.
	sqlQueryBase_SelectFileEvents :=
		`(SELECT
	tmp.local_day,
	count(tmp.local_day),
	max(tmp.local_time) AS local_time,
	3 AS event_type_id
FROM
	(SELECT
		event_time at time zone $1 AS local_time,
		date(event_time at time zone $2) AS local_day
	FROM public.download_events
	WHERE file_id = $3
	ORDER BY local_time ASC) AS tmp
GROUP BY tmp.local_day
ORDER BY tmp.local_day)

UNION

(SELECT
	NULL,
	NULL,
	event_time at time zone $4 AS local_time,
	event_type_id
FROM public.simple_events
WHERE file_id = $5
ORDER BY event_time ASC)

ORDER BY local_time DESC`

	s.sqlQueryTexts[SqlQueryIndexForQueryToSelectFileEventsAll] = sqlQueryBase_SelectFileEvents +
		sqlQueryTerminator

	s.sqlQueryTexts[SqlQueryIndexForQueryToSelectFileEventsLimited] = sqlQueryBase_SelectFileEvents + ` LIMIT $6` +
		sqlQueryTerminator
}

// Open начинает работу со хранилищем.
func (s *Storage) Open() {
	s.start()
}

// start запускает хранилище.
func (s *Storage) start() {
	s.connectorWG.Add(1)

	go s.keepConnected()
}

// keepConnected держит соединение со хранилищем.
func (s *Storage) keepConnected() {
	defer s.connectorWG.Done()

	s.logger.Info().Msg(msgStorageConnectorStart)

	var err error

	for {
		if s.isShutdownRequired() {
			break
		}

		// Если нет соединения со хранилищем, то устанавливаем его.
		// Если соединение есть, но не работает, переустанавливаем его.
		// Если соединение есть и работает, то засыпаем.
		if !s.isConnectionNull() {
			err = s.Ping()
			if err == nil {
				// Всё хорошо.

				// Если ранее база данных была недоступна, сообщаем о доступности.
				if !s.IsReady() {
					s.logger.Info().Msg(msgStorageConnected)
				}

				s.setReadinessState(true)
				time.Sleep(time.Second * config.StorageKeeperReadinessCheckIntervalSec)
				continue
			}
		}

		// Что-то пошло не так.
		s.setReadinessState(false)

		// Если не прошёл Ping или упало соединение, то пишем в журнал причину.
		if s.isConnectionNull() {
			s.logger.Info().Msg(msgStorageNotConnected)
		} else {
			if err != nil {
				s.logger.Err(err).Msg(msgStoragePingFailure)
			}
		}

		// Устанавливаем соединение с базой данных.
		s.logger.Info().Msg(msgStorageConnecting)
		err = s.connect()
		if err != nil {
			s.prometheusMetrics.DatabaseConnectionErrorsCount.Inc()
			s.logger.Err(err).Msg(msgStorageConnectionFailure)

			// Чтобы не загадить журнал ошибками, ждём немного.
			time.Sleep(time.Second * config.StorageConnectorDelaySecAfterError)
			continue
		}

		// Подготавливаем statements.
		err = s.prepareSqlQueryStatements()
		if err != nil {
			s.logger.Err(err).Msg(msgStorageConnectionFailure)

			// Чтобы не загадить журнал ошибками, ждём немного.
			time.Sleep(time.Second * config.StorageConnectorDelaySecAfterError)
			continue
		}

		s.logger.Info().Msg(msgStoragePreparedStatementsReady)

		s.isReady.Store(true)
	}

	s.logger.Info().Msg(msgStorageConnectorStop)
}

// connect соединяется с сервером базы данных и проверяет соединение.
func (s *Storage) connect() (err error) {
	const (
		// SqlDriverName -- название SQL драйвера.
		SqlDriverName = "pgx"

		// PostgreSqlDsnPrefix -- префикс DNS для PostgreSQL.
		PostgreSqlDsnPrefix = "postgresql"
	)

	var userWithPassword string
	if len(s.config.PostgrePassword) < 1 {
		userWithPassword = s.config.PostgreUser
	} else {
		userWithPassword = fmt.Sprintf(
			"%s:%s",
			s.config.PostgreUser,
			s.config.PostgrePassword,
		)
	}

	dsn := fmt.Sprintf(
		"%s://%s@%s:%d/%s?%s",
		PostgreSqlDsnPrefix,
		userWithPassword,
		s.config.PostgreHost,
		s.config.PostgrePort,
		s.config.PostgreDatabase,
		s.config.PostgreParameters,
	)

	s.connection, err = sqlx.Open(SqlDriverName, dsn)
	if err != nil {
		return err
	}

	err = s.Ping()
	if err != nil {
		return err
	}

	s.logger.Info().Msg(msgStorageConnected)

	return nil
}

// Ping производит проверку связи.
func (s *Storage) Ping() (err error) {
	return s.connection.Ping()
}

// IsReady возвращает состояние готовности.
func (s *Storage) IsReady() bool {
	return s.isReady.Load()
}

// setReadinessState устанавливает новое состояние готовности.
func (s *Storage) setReadinessState(newState bool) {
	s.isReady.Store(newState)
}

// isConnectionNull проверяет соединение на равенство пустому указателю.
func (s *Storage) isConnectionNull() bool {
	return s.connection == nil
}

// isShutdownRequired говорит, нужно ли останавливать хранилище.
// Если компонент, выполнивший этот метод, получил ответ 'true', то он теряет
// право опрашивать этот метод, поскольку, важно, чтобы все компоненты
// хранилища смогли получить сигнал завершения.
func (s *Storage) isShutdownRequired() bool {
	select {
	case <-s.close:
		return true
	default:
		return false
	}
}

// Wait ждёт готовности хранилища.
func (s *Storage) Wait() {
	for {
		if s.isReady.Load() {
			break
		} else {
			time.Sleep(time.Millisecond * config.StorageReadinessWaitIntervalMs)
		}
	}
}

// Close останавливает хранилище.
func (s *Storage) Close() (err error) {
	// Говорим всем компонентам, что нужно остановиться.
	var i byte
	for i = 1; i <= s.componentsCount; i++ {
		s.close <- true
	}

	// Ждём завершения всех компонентов.
	s.connectorWG.Wait()

	// Разрываем соединение с базой данных.
	err = s.disconnect()
	if err != nil {
		return err
	}

	return nil
}

// disconnect отсоединяется от сервера базы данных.
func (s *Storage) disconnect() (err error) {
	if s.isConnectionNull() {
		return errors.New(ErrStorageConnectionIsNull)
	}

	err = s.connection.Close()
	if err != nil {
		return err
	}

	return nil
}
