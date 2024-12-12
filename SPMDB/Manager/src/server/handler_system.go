package server

import (
	"net/http"
	"spdb_manager/cli/printer"

	"github.com/julienschmidt/httprouter"
)

func (srv *Server) hSystemShutdown(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	printer.PrintInfo("Server is shutting down ...")
	go srv.stop()
}
