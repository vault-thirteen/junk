package httpserver

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	iHttp "github.com/vault-thirteen/junk/gfe/internal/http"
)

// ErrorMsgAccessTokenCheck -- сообщение об ошибке проверки токена доступа.
const ErrorMsgAccessTokenCheck = "access token check"

// middlewareAuthentication -- HTTP middleware для проверки аутентификации.
func (hs *HttpServer) middlewareAuthentication(nextHandler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Счётчик запросов.
		hs.prometheusMetrics.RequestsCount.With(prometheus.Labels{
			config.MetricParameterPath: r.URL.Path,
		}).Inc()

		// Проверка веб-токена.
		var err = hs.checkAccessToken(r)
		if err != nil {
			hs.logger.Err(err).Msg(ErrorMsgAccessTokenCheck)
			iHttp.RespondWithNotAuthorizedError(w)
			return
		}

		nextHandler(w, r, ps)
	}
}
