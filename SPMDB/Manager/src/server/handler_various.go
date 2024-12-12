package server

import (
	"fmt"
	"net/http"
	"spdb_manager/settings"

	"github.com/julienschmidt/httprouter"
)

const RootPageContents = `<html>
	<head>
		<meta http-equiv="Refresh" content="3; url=/` + settings.GuiUrlFolder + `/">
	</head>
	<body>
		Redirecting ...
	</body>
</html>`

func (srv *Server) hIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	_, err := fmt.Fprint(w, RootPageContents)
	srv.logError(err)
}
