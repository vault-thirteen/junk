package http

import (
	"net/http"

	"github.com/rs/zerolog"
)

func RespondWithNotReadyStatus(logger *zerolog.Logger, w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusServiceUnavailable)

	_, err := w.Write([]byte(reason))
	logOnError(logger, err)
}

func logOnError(logger *zerolog.Logger, err error) {
	if err != nil {
		logger.Err(err).Send()
	}
}
