package hs

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	iHttp "github.com/vault-thirteen/junk/SSE2/internal/http"
)

func (hs *HttpServer) handlerLiveness(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
}

func (hs *HttpServer) handlerReadiness(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, err := hs.service.GetReadinessState()
	if err != nil {
		iHttp.RespondWithNotReadyStatus(hs.logger, w, err.Error())
		return
	}

	return
}
