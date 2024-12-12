package server

import "net/http"

func (srv *Server) hBadRequest(w http.ResponseWriter, err error) {
	srv.logError(err)
	w.WriteHeader(http.StatusBadRequest)
}

func (srv *Server) hInternalServerError(w http.ResponseWriter, err error) {
	srv.logError(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (srv *Server) hNotFound(w http.ResponseWriter, err error) {
	srv.logError(err)
	w.WriteHeader(http.StatusNotFound)
}
