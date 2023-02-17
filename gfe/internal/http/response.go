package http

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

// RespondWithNotReadyStatus отвечает на HTTP запрос сообщением о недоступности
// сервиса.
func RespondWithNotReadyStatus(logger *zerolog.Logger, w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusServiceUnavailable)

	_, err := w.Write([]byte(reason))
	logOnError(logger, err)
}

// RespondWithNotAuthorizedError отвечает на HTTP запрос ошибкой типа
// "не авторизован".
func RespondWithNotAuthorizedError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// RespondWithBadRequestError отвечает на HTTP запрос ошибкой типа
// "неправильный запрос".
func RespondWithBadRequestError(w http.ResponseWriter, reason string) {
	http.Error(w, reason, http.StatusBadRequest)
}

// RespondWithInternalServerError отвечает на HTTP запрос ошибкой типа
// "внутренняя ошибка сервера".
func RespondWithInternalServerError(w http.ResponseWriter, reason string) {
	http.Error(w, reason, http.StatusInternalServerError)
}

// RespondWithJsonObject отвечает на HTTP запрос JSON объектом.
func RespondWithJsonObject(logger *zerolog.Logger, w http.ResponseWriter, object interface{}) {
	buf, err := json.Marshal(object)
	if err != nil {
		RespondWithInternalServerError(w, err.Error())
		return
	}

	w.Header().Set(HttpHeaderContentType, MimeTypeApplicationJson)

	_, err = w.Write(buf)
	logOnError(logger, err)
}

// logOnError -- пишет ошибку в журнал, если она есть.
func logOnError(logger *zerolog.Logger, err error) {
	if err != nil {
		logger.Err(err).Send()
	}
}
