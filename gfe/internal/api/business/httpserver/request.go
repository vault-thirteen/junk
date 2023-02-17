package httpserver

import (
	"net/http"
	"time"

	iHttp "github.com/vault-thirteen/junk/gfe/internal/http"
	"github.com/vault-thirteen/junk/gfe/internal/message"
	"github.com/vault-thirteen/junk/gfe/pkg/models/file"
	fer "github.com/vault-thirteen/junk/gfe/pkg/models/fileeventsrequest"
)

// getRequestToGetFileEvents читает параметры HTTP запроса для выдачи истории
// по файлу.
func (hs *HttpServer) getRequestToGetFileEvents(r *http.Request) (req *fer.FileEventsRequest, err error) {
	req = new(fer.FileEventsRequest)

	// Проверяем часовой пояс клиента.
	req.ClientTimeZone = r.Header.Get(iHttp.HttpHeaderClientTimeZoneName)
	if len(req.ClientTimeZone) < 1 {
		return nil, message.ErrClientTimeZoneNotSet
	}

	_, err = time.LoadLocation(req.ClientTimeZone)
	if err != nil {
		return nil, err
	}

	// Проверяем идентификатор файла.
	err = r.ParseForm()
	if err != nil {
		return nil, err
	}

	var buf = r.Form.Get(iHttp.HttpQueryParameterFileID)
	req.FileID = file.NewID(buf)
	if len(req.FileID) < 1 {
		return nil, message.ErrFileIDNotSet
	}

	return req, nil
}
