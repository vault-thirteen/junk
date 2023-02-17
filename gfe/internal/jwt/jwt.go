package jwt

import (
	"crypto/rsa"
	"os"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	"github.com/vault-thirteen/junk/gfe/internal/cypher"
	"github.com/vault-thirteen/junk/gfe/internal/dsn"
	"github.com/vault-thirteen/junk/gfe/internal/message"
	"github.com/vault-thirteen/junk/gfe/pkg/models/keysourcetype"
)

// MsgFDebugConfig -- формат сообщения для отладки настроек.
const MsgFDebugConfig = "jwt configuration: %+v"

// Jwt -- JWT инфраструктура.
type Jwt struct {
	// Настройки JWT.
	config *config.Jwt

	// Публичный RSA ключ для проверки JWT токенов.
	rsaPublicKey *rsa.PublicKey
}

// NewJwt -- конструктор JWT инфраструктуры.
func NewJwt(logger *zerolog.Logger) (*Jwt, error) {
	j := new(Jwt)

	err := j.init()
	if err != nil {
		return nil, err
	}

	logger.Debug().Msg(pretty.Sprintf(MsgFDebugConfig, j.config))

	return j, nil
}

// init -- инициализация инфраструктуры JWT.
func (j *Jwt) init() (err error) {
	j.config, err = config.GetJwtConfig()
	if err != nil {
		return err
	}

	err = j.readJwtPublicKey()
	if err != nil {
		return err
	}

	return nil
}

// readJwtPublicKey загружает публичный ключ для проверки JWT токенов.
func (j *Jwt) readJwtPublicKey() (err error) {
	// Читаем публичный RSA ключ в формате PEM в переменную.
	var publicKeyPemText string

	switch j.config.JwtKeySourceType {
	case keysourcetype.EnvironmentVariable:
		// Берём ключ из переменной окружения.
		publicKeyPemText = config.UnScreenJwtPublicKeyString(j.config.JwtKeyValue)

	case keysourcetype.File:
		// Берём ключ из файла.
		var filePath string
		filePath, err = dsn.GetFilePathFromDsn(j.config.JwtKeyDsn)
		if err != nil {
			return err
		}

		publicKeyPemText, err = j.getTokenKeyFromFile(filePath)
		if err != nil {
			return err
		}

	case keysourcetype.Vault:
		// Берём ключ из Vault.
		//TODO: сделать получение ключа из Vault.
		// Информация по настройкам Vault пока не доступна.
		return errors.New("vault is not accessible")

	default:
		return errors.Errorf(message.ErrFJwtKeySourceTypeUnsupported, j.config.JwtKeySourceType)
	}

	// Распознаём и сохраняем ключ для быстрого доступа к нему.
	j.rsaPublicKey, err = cypher.ParseRsaPublicKey(publicKeyPemText)
	if err != nil {
		return errors.Wrap(err, "ParseRsaPublicKey")
	}

	return nil
}

// getTokenKeyFromFile читает ключ JWT токена из файла.
func (j *Jwt) getTokenKeyFromFile(filePath string) (tokenKeyText string, err error) {
	var buf []byte
	buf, err = os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// GetRsaPublicKey возвращает публичный RSA ключ для проверки JWT токенов.
func (j *Jwt) GetRsaPublicKey() *rsa.PublicKey {
	return j.rsaPublicKey
}
