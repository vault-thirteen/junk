package message

import "github.com/pkg/errors"

// Форматы ошибок.
const (
	// ErrFJwtKeySourceTypeUnsupported -- сообщение об ошибке "неподдерживаемый тип источника JWT ключа".
	ErrFJwtKeySourceTypeUnsupported = "unsupported jwt key source type: %v"
)

// Ошибки.
var (
	// ErrClientTimeZoneNotSet -- сообщение об ошибке "часовой пояс клиента не установлен".
	ErrClientTimeZoneNotSet = errors.New("client time zone is not set")

	// ErrFileIDNotSet -- сообщение об ошибке "идентификатор файла не установлен".
	ErrFileIDNotSet = errors.New("file id is not set")

	// ErrAuthorizationHeaderFormatUnsupported -- сообщение об ошибке "формат заголовка авторизации не поддерживается".
	ErrAuthorizationHeaderFormatUnsupported = errors.New("authorization header format is not supported")

	// ErrAccessTokenNotValid -- сообщение об ошибке "токен доступа недействителен".
	ErrAccessTokenNotValid = errors.New("access token is not valid")

	// ErrIsAlreadyStarted -- сообщение об ошибке "уже запущен".
	ErrIsAlreadyStarted = errors.New("is already started")

	// ErrIsNotStarted -- сообщение об ошибке "не запущен".
	ErrIsNotStarted = errors.New("is not started")
)
