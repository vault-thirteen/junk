package httpserver

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	iHttp "github.com/vault-thirteen/junk/gfe/internal/http"
	"github.com/vault-thirteen/junk/gfe/internal/message"
)

// AuthorizationHeaderPartsCountExpected -- ожидаемое количество частей
// заголовка авторизации.
const AuthorizationHeaderPartsCountExpected = 2

// checkAccessToken проверяет действительность токена доступа.
// Токен доступа использует стандарт JWT (веб токен JSON).
// Веб токен должен использовать RSA шифрование.
// Публичный RSA ключ достаём из внешнего хранилища ключей.
func (hs *HttpServer) checkAccessToken(r *http.Request) (err error) {
	var token string
	token, err = hs.readAuthorizationToken(r)
	if err != nil {
		return err
	}

	var tokenObject *jwt.Token
	tokenObject, err = jwt.Parse(token, hs.accessTokenParser)
	if err != nil {
		return err
	}

	if !tokenObject.Valid {
		return message.ErrAccessTokenNotValid
	}

	return nil
}

// readAuthorizationToken читает токен доступа (аутентификации, авторизации) из
// входящего HTTP запроса.
func (hs *HttpServer) readAuthorizationToken(r *http.Request) (token string, err error) {
	buffer := r.Header.Get(iHttp.HttpHeaderAuthorization)

	tokenParts := strings.Split(buffer, " ")
	if len(tokenParts) != AuthorizationHeaderPartsCountExpected {
		return "", message.ErrAuthorizationHeaderFormatUnsupported
	}

	if tokenParts[0] != iHttp.HttpHeaderAuthorizationTypeBearer {
		return "", message.ErrAuthorizationHeaderFormatUnsupported
	}

	return tokenParts[1], nil
}

// Распознаватель токена доступа. Библиотека JWT использует этот метод для
// получения ключа для проверки подписи веб токена.
// Параметр 'token' нужен для соответствия типу 'jwt.Keyfunc'.
func (hs *HttpServer) accessTokenParser(_ *jwt.Token) (key interface{}, err error) {
	return hs.jwtRsaPublicKey, nil
}
