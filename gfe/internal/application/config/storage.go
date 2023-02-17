package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Различные настройки.
const (
	// StorageConnectorDelaySecAfterError -- период времени (в секундах), в
	// течение которого коннектор базы данных делает паузу после неудачной
	// попытки соединения с сервером базы данных.
	StorageConnectorDelaySecAfterError = 5

	// StorageReadinessWaitIntervalMs -- период времени (в миллисекундах), в
	// течение которого происходит задержка перед следующей попыткой проверки
	// готовности хранилища (базы данных).
	StorageReadinessWaitIntervalMs = 100

	// StorageKeeperReadinessCheckIntervalSec -- период времени (в секундах), в
	// течение которого коннектор базы данных делает паузу после удачного пинга
	// сервера базы данных.
	StorageKeeperReadinessCheckIntervalSec = 5

	// StorageQueryTimeoutSec -- предельный период времени (в секундах), в
	// течение которого запрос в базу данных должен выполниться.
	StorageQueryTimeoutSec = 5 * 60
)

// Storage -- настройка хранилища (базы данных).
type Storage struct {
	// Хост.
	PostgreHost string `split_words:"true" default:"localhost"`

	// Порт.
	PostgrePort uint `split_words:"true" default:"5432"`

	// Пользователь.
	PostgreUser string `split_words:"true"`

	// Пароль.
	PostgrePassword string `split_words:"true"`

	// База данных
	PostgreDatabase string `split_words:"true"`

	// Строка дополнительных параметров подключения.
	// Не содержит символа '?'.
	PostgreParameters string `split_words:"true"`
}

// NewStorage -- создатель настроек хранилища.
func NewStorage(envPrefix string) (cfg *Storage, err error) {
	cfg = new(Storage)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsValid проверяет годность настройки хранилища.
func (s *Storage) IsValid() (bool, error) {
	if len(s.PostgreHost) < 1 {
		return false, ErrHost
	}

	if s.PostgrePort < 1 {
		return false, ErrPort
	}

	if len(s.PostgreUser) < 1 {
		return false, ErrUser
	}

	if len(s.PostgreDatabase) < 1 {
		return false, ErrDatabase
	}

	return true, nil
}

// GetStorageConfig получает настройки хранилища.
func GetStorageConfig() (storageConfig *Storage, err error) {
	storageConfig, err = NewStorage(environmentVariablePrefixApplication)
	if err != nil {
		return nil, err
	}

	_, err = storageConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return storageConfig, nil
}
