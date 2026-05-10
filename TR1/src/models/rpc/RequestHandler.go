package rm

import (
	"net/http"
)

type RequestHandler = func(ar *Request, hr *http.Request, hrw http.ResponseWriter)
