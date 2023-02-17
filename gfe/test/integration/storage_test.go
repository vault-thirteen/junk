package integration

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	"github.com/vault-thirteen/junk/gfe/pkg/models/downloadevent"
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/simpleevent"
)

// Вспомогательный объект для проведения интеграционного тестирования хранилища.
type Storage struct {
	Settings   *config.Storage
	Connection *sqlx.DB
}

// Данные из базы данных.
type StorageData struct {
	EventTypes     []event.Type
	SimpleEvents   []simpleevent.SimpleEvent
	DownloadEvents []downloadevent.DownloadEvent
}

// NewStorage -- конструктор объекта.
func NewStorage() (s *Storage, err error) {
	s = new(Storage)

	// Настройки.
	s.Settings = &config.Storage{
		PostgreHost:       "localhost",
		PostgrePort:       5432,
		PostgreUser:       IntegrationTestIdPrefix,
		PostgrePassword:   IntegrationTestIdPrefix,
		PostgreDatabase:   IntegrationTestIdPrefix,
		PostgreParameters: "",
	}

	// Connection.
	err = s.connect()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) connect() (err error) {
	const (
		// SqlDriverName -- название SQL драйвера.
		SqlDriverName = "pgx"

		// PostgreSqlDsnPrefix -- префикс DNS для PostgreSQL.
		PostgreSqlDsnPrefix = "postgresql"
	)

	var userWithPassword string
	if len(s.Settings.PostgrePassword) < 1 {
		userWithPassword = s.Settings.PostgreUser
	} else {
		userWithPassword = fmt.Sprintf(
			"%s:%s",
			s.Settings.PostgreUser,
			s.Settings.PostgrePassword,
		)
	}

	dsn := fmt.Sprintf(
		"%s://%s@%s:%d/%s?%s",
		PostgreSqlDsnPrefix,
		userWithPassword,
		s.Settings.PostgreHost,
		s.Settings.PostgrePort,
		s.Settings.PostgreDatabase,
		s.Settings.PostgreParameters,
	)

	s.Connection, err = sqlx.Open(SqlDriverName, dsn)
	if err != nil {
		return err
	}

	err = s.Connection.Ping()
	if err != nil {
		return err
	}

	return nil
}

// Stop останавливает объект.
func (s *Storage) Stop() (err error) {
	err = s.cleanDatabase()
	if err != nil {
		return err
	}

	err = s.disconnect()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) cleanDatabase() (err error) {
	queries := []string{
		`TRUNCATE TABLE simple_events;`,
		`TRUNCATE TABLE download_events;`,
	}

	for _, query := range queries {
		_, err = s.Connection.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) disconnect() (err error) {
	return s.Connection.Close()
}

// ReadData читает содержимое базы данных.
func (s *Storage) ReadData() (data *StorageData, err error) {
	data = &StorageData{}

	// Типы событий.
	query := `SELECT id, description_ru, description_en, is_simple, is_aggregated
FROM public.event_types;`

	err = s.Connection.Select(&data.EventTypes, query)
	if err != nil {
		return nil, err
	}

	// Простые события.
	query = `SELECT id, user_id, file_id, event_type_id, event_time
FROM public.simple_events;`

	err = s.Connection.Select(&data.SimpleEvents, query)
	if err != nil {
		return nil, err
	}

	// События типа 'скачивание'.
	query = `SELECT id, user_id, file_id, event_time
FROM public.download_events;`

	err = s.Connection.Select(&data.DownloadEvents, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}
