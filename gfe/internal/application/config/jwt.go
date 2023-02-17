package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/vault-thirteen/junk/gfe/pkg/models/keysourcetype"
)

// JwtSymbolForCrLf -- символ для экранирования CR/LF при передаче ключа через
// переменную окружения.
const JwtSymbolForCrLf = "."

// Jwt -- настройка JWT инфраструктуры.
type Jwt struct {
	// Тип источника публичного RSA ключа для проверки подписи JWT токена.
	JwtKeySourceType keysourcetype.KeySourceType `split_words:"true"`

	// Строка, уточняющая способ получения публичного RSA ключа. Хранит
	// параметры соединения с Vault или другими источниками. Если тип источника
	// ключа -- переменная окружения, то этот параметр пуст.
	JwtKeyDsn string `split_words:"true"`

	// Значение публичного RSA ключа для проверки подписи JWT токена, если тип
	// источника ключа -- переменная окружения, в иных случаях -- пустота. Если
	// значение не пусто, то должно быть в формате PEM.
	JwtKeyValue string `split_words:"true"`
}

// NewJwt -- создатель настроек JWT инфраструктуры.
func NewJwt(envPrefix string) (cfg *Jwt, err error) {
	cfg = new(Jwt)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsValid проверяет годность настройки HTTP сервера.
func (c *Jwt) IsValid() (bool, error) {
	if !c.JwtKeySourceType.IsValid() {
		return false, errors.Errorf(ErrFKeySourceTypeNotValid, c.JwtKeySourceType)
	}

	// Если используется переменная окружения, то DSN должен быть пуст, а
	// значение ключа должно быть заполнено.
	if c.JwtKeySourceType.IsEnvironmentVariable() {
		if len(c.JwtKeyDsn) > 0 {
			return false, ErrKeyDsnIsSetButNotUsed
		}
		if len(c.JwtKeyValue) < 1 {
			return false, ErrKeyValueIsNotSetButUsed
		}
	}

	// Если используется файл, то DSN должен быть задан, а значение ключа
	// должно быть пусто.
	if c.JwtKeySourceType.IsFile() {
		if len(c.JwtKeyDsn) < 1 {
			return false, ErrKeyDsnIsNotSetButUsed
		}
		if len(c.JwtKeyValue) > 0 {
			return false, ErrKeyValueIsSetButNotUsed
		}
	}

	// Если используется Vault, то DSN должен быть задан, а значение ключа
	// должно быть пусто.
	if c.JwtKeySourceType.IsVault() {
		if len(c.JwtKeyDsn) < 1 {
			return false, ErrKeyDsnIsNotSetButUsed
		}
		if len(c.JwtKeyValue) > 0 {
			return false, ErrKeyValueIsSetButNotUsed
		}
	}

	return true, nil
}

// GetJwtConfig получает настройки JWT инфраструктуры.
func GetJwtConfig() (jwtConfig *Jwt, err error) {
	jwtConfig, err = NewJwt(environmentVariablePrefixApplication)
	if err != nil {
		return nil, err
	}

	_, err = jwtConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return jwtConfig, nil
}

// UnScreenJwtPublicKeyString де-экранирует CR/LF в публичном ключе при передаче
// ключа через переменную окружения.
func UnScreenJwtPublicKeyString(jwtKeyPemTextScreened string) (jwtKeyPemText string) {
	// '.' -> CR/LF.
	return strings.ReplaceAll(jwtKeyPemTextScreened, JwtSymbolForCrLf, "\r\n")
}

// ScreenJwtPublicKeyString экранирует CR/LF в публичном ключе для передачи
// ключа через переменную окружения.
func ScreenJwtPublicKeyString(jwtKeyPemText string) (jwtKeyPemTextScreened string) {
	// Поскольку стиль переноса строки может быть разным, применяем хитрость.
	// Сначала экранируем самый сложный формат. Если останутся другие форматы,
	// то экранируем и их.

	// CR/LF -> '.'.
	buffer := strings.ReplaceAll(jwtKeyPemText, "\r\n", JwtSymbolForCrLf)

	// CR -> '.'.
	buffer = strings.ReplaceAll(buffer, "\r", JwtSymbolForCrLf)

	// LF -> '.'.
	buffer = strings.ReplaceAll(buffer, "\n", JwtSymbolForCrLf)

	return buffer
}
