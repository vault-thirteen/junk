package cypher

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

//TODO: Доделать этот пакет, когда информация о Vault станет доступной.

const (
	// RsaKeyTypeRsaPublicKey -- тип RSA ключа :: RSA PUBLIC KEY.
	RsaKeyTypeRsaPublicKey = "RSA PUBLIC KEY"

	// RsaKeyTypePublicKey -- тип RSA ключа :: PUBLIC KEY.
	RsaKeyTypePublicKey = "PUBLIC KEY"
)

// Ошибки.
var (
	// ErrRsaKeyNotFound -- ошибка "RSA ключ не найден".
	ErrRsaKeyNotFound = errors.New("rsa key is not found")

	// ErrRsaKeyTypeUnsupported -- ошибка "неподдерживаемый тип RSA ключа".
	ErrRsaKeyTypeUnsupported = errors.New("unsupported rsa key type")

	// ErrRsaKeyIsNotPublic -- ошибка "RSA ключ -- не публичный".
	ErrRsaKeyIsNotPublic = errors.New("rsa key is not public")
)

// ParseRsaPublicKey получает из публичного RSA ключа в формате PEM ключ в виде
// объекта.
func ParseRsaPublicKey(publicKeyPEM string) (rsaPublicKey *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, ErrRsaKeyNotFound
	}

	switch block.Type {
	case RsaKeyTypePublicKey:
		rsaPublicKey, err = parseRsaPublicKeyPKIX(block.Bytes)
		if err != nil {
			return nil, err
		}

	default:
		return nil, ErrRsaKeyTypeUnsupported
	}

	return rsaPublicKey, nil
}

// parseRsaPublicKeyPKIX -- распознаёт публичный RSA ключ для ключа в PEM формате
// "PUBLIC KEY".
func parseRsaPublicKeyPKIX(blockBytes []byte) (rsaPublicKey *rsa.PublicKey, err error) {
	var genericKey interface{}
	genericKey, err = x509.ParsePKIXPublicKey(blockBytes)
	if err != nil {
		return nil, err
	}

	var ok bool
	rsaPublicKey, ok = genericKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrRsaKeyIsNotPublic
	}

	return rsaPublicKey, nil
}
