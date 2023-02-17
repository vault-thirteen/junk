package config

import "github.com/pkg/errors"

// Форматы сообщений об ошибке.
const (
	// ErrFKeySourceTypeNotValid -- ошибка "тип источника JWT ключа не пригоден".
	ErrFKeySourceTypeNotValid = "jwt key source type '%v' is not valid"
)

// Сообщения об ошибках.
var (
	// ErrHost -- ошибка "хост пуст".
	ErrHost = errors.New("host is empty")

	// ErrPort -- ошибка "порт пуст".
	ErrPort = errors.New("port is empty")

	// ErrConsumerGroupID -- ошибка "идентификатор группы потребителя пуст".
	ErrConsumerGroupID = errors.New("consumer group id is empty")

	// ErrBrokerAddressListEmpty -- ошибка "список адресов брокеров пуст".
	ErrBrokerAddressListEmpty = errors.New("broker address list is empty")

	// ErrTopicListEmpty -- ошибка "список тем пуст".
	ErrTopicListEmpty = errors.New("topic list is empty")

	// ErrKeyDsnIsSetButNotUsed -- ошибка "DSN JWT ключа задан, но не использован".
	ErrKeyDsnIsSetButNotUsed = errors.New("jwt key dsn is set but not used")

	// ErrKeyDsnIsNotSetButUsed -- ошибка "DSN JWT ключа не задан, но использован".
	ErrKeyDsnIsNotSetButUsed = errors.New("jwt key dsn is not set but used")

	// ErrKeyValueIsNotSetButUsed -- ошибка "значение JWT ключа не задано, но использовано".
	ErrKeyValueIsNotSetButUsed = errors.New("jwt key value is not set but used")

	// ErrKeyValueIsSetButNotUsed -- ошибка "значение JWT ключа задано, но не использовано".
	ErrKeyValueIsSetButNotUsed = errors.New("jwt key value is set but not used")

	// ErrUser -- ошибка "пользователь пуст".
	ErrUser = errors.New("user is empty")

	// ErrDatabase -- ошибка "база данных пуст".
	ErrDatabase = errors.New("database is empty")
)
