package integration

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/pkg/errors"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	iHttp "github.com/vault-thirteen/junk/gfe/internal/http"
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/history"
	"go.uber.org/multierr"
)

const (
	// ProtocolSchemaFull -- полная схема протокола.
	ProtocolSchemaFull = "http://"

	// TimeZoneForTests -- название часового пояса клиента, проводящего
	// проверку.
	TimeZoneForTests = "Africa/Johannesburg" // UTC + 02:00, no DST.
)

// ErrFHttpStatusCodeWrong -- ошибка "http статус код неправилен".
const ErrFHttpStatusCodeWrong = "http status code is wrong, expected=%v, actual=%v"

// Вспомогательный объект для проведения интеграционного тестирования API.
type API struct {
	// Родительский объект.
	parent *Test

	// Настройки системного HTTP сервера.
	SystemHttpServerSettings *config.HttpServer

	// Настройки HTTP сервера бизнес логики.
	BusinessLogicsHttpServerSettings *config.HttpServer
}

// Данные, полученные из API в ходе проверки.
type APIData struct {
	// 1. Данные с системного сервера.

	// 1.1. Статус, полученный от хендлера готовности.
	ReadinessStatus int

	// 1.2. Статус, полученный от хендлера метрик.
	MetricsStatus int

	// 2. Данные с сервера бизнес логики.

	// 2.1. Статус, полученный от хендлера жизни сервиса.
	LivenessStatus int

	// 2.2. Данные, полученные от хендлера типов событий.
	FileEventTypes *FileEventTypesData

	// 2.3. Данные, полученные от хендлера всех событий.
	FileEventsAll *FileEventsAllData

	// 2.4. Данные, полученные от хендлера недавних событий.
	FileEventsLastN *FileEventsLastNData
}

// Данные, полученные от хендлера типов событий.
type FileEventTypesData struct {
	Data []event.Type
}

// Данные, полученные от хендлера всех событий.
type FileEventsAllData struct {
	Data history.History
}

// Данные, полученные от хендлера недавних событий.
type FileEventsLastNData struct {
	Data history.History
}

// NewAPI -- конструктор объекта.
func NewAPI(parent *Test) (a *API, err error) {
	a = new(API)

	a.parent = parent

	// Настройки системного HTTP сервера (метрик).
	a.SystemHttpServerSettings = &config.HttpServer{
		HttpServerHost: "localhost",
		HttpServerPort: 2001,
	}

	// Настройки HTTP сервера бизнес логики.
	a.BusinessLogicsHttpServerSettings = &config.HttpServer{
		HttpServerHost: "localhost",
		HttpServerPort: 2002,
	}

	return a, nil
}

// CheckData делает запросы к HTTP серверам сервиса и отдаёт их ответы.
func (a *API) ReadData() (data *APIData, err error) {
	data = &APIData{}

	data.ReadinessStatus, err = a.readReadinessStatus()
	if err != nil {
		return nil, err
	}

	data.MetricsStatus, err = a.readMetricsStatus()
	if err != nil {
		return nil, err
	}

	data.LivenessStatus, err = a.readLivenessStatus()
	if err != nil {
		return nil, err
	}

	data.FileEventTypes, err = a.readEventTypes()
	if err != nil {
		return nil, err
	}

	data.FileEventsAll, err = a.readAllEvents()
	if err != nil {
		return nil, err
	}

	data.FileEventsLastN, err = a.readLastEvents()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// readReadinessStatus читает статус готовоности сервиса.
func (a *API) readReadinessStatus() (readinessStatus int, err error) {
	serverAddress := net.JoinHostPort(
		a.SystemHttpServerSettings.HttpServerHost,
		strconv.Itoa(int(a.SystemHttpServerSettings.HttpServerPort)),
	)

	requestPath := ProtocolSchemaFull + path.Join(serverAddress, "ready")

	var httpResponse *http.Response
	httpResponse, err = http.Get(requestPath)
	if err != nil {
		return 0, err
	}

	return httpResponse.StatusCode, nil
}

// readMetricsStatus читает статус доступности метрик сервиса.
func (a *API) readMetricsStatus() (metricsStatus int, err error) {
	serverAddress := net.JoinHostPort(
		a.SystemHttpServerSettings.HttpServerHost,
		strconv.Itoa(int(a.SystemHttpServerSettings.HttpServerPort)),
	)

	requestPath := ProtocolSchemaFull + path.Join(serverAddress, "metrics")

	var httpResponse *http.Response
	httpResponse, err = http.Get(requestPath)
	if err != nil {
		return 0, err
	}

	return httpResponse.StatusCode, nil
}

// readLivenessStatus читает статус жизни сервиса.
func (a *API) readLivenessStatus() (livenessStatus int, err error) {
	serverAddress := net.JoinHostPort(
		a.BusinessLogicsHttpServerSettings.HttpServerHost,
		strconv.Itoa(int(a.BusinessLogicsHttpServerSettings.HttpServerPort)),
	)

	requestPath := ProtocolSchemaFull + path.Join(serverAddress, "live")

	var httpResponse *http.Response
	httpResponse, err = http.Get(requestPath)
	if err != nil {
		return 0, err
	}

	return httpResponse.StatusCode, nil
}

// readEventTypes читает типы событий сервиса.
func (a *API) readEventTypes() (data *FileEventTypesData, err error) {
	data = &FileEventTypesData{}

	serverAddress := net.JoinHostPort(
		a.BusinessLogicsHttpServerSettings.HttpServerHost,
		strconv.Itoa(int(a.BusinessLogicsHttpServerSettings.HttpServerPort)),
	)

	requestPath := ProtocolSchemaFull + path.Join(serverAddress, "file-event", "types")

	httpRequestHeaders := make(http.Header)
	httpRequestHeaders.Add(a.parent.makeAuthorizationHeader())

	var httpRequestUrl *url.URL
	httpRequestUrl, err = url.Parse(requestPath)
	if err != nil {
		return nil, err
	}

	httpRequest := &http.Request{
		URL:    httpRequestUrl,
		Method: http.MethodGet,
		Header: httpRequestHeaders,
	}

	var httpResponse *http.Response
	httpResponse, err = http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.Errorf(ErrFHttpStatusCodeWrong, http.StatusOK, httpResponse.StatusCode)
	}

	var httpResponseBody []byte
	httpResponseBody, err = io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		derr := httpResponse.Body.Close()
		err = multierr.Combine(err, derr)
	}()

	err = json.Unmarshal(httpResponseBody, &data.Data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// readAllEvents читает все события файла.
func (a *API) readAllEvents() (data *FileEventsAllData, err error) {
	data = &FileEventsAllData{}

	serverAddress := net.JoinHostPort(
		a.BusinessLogicsHttpServerSettings.HttpServerHost,
		strconv.Itoa(int(a.BusinessLogicsHttpServerSettings.HttpServerPort)),
	)

	requestPath := ProtocolSchemaFull + path.Join(serverAddress, "file-events", "all")

	httpRequestHeaders := make(http.Header)
	httpRequestHeaders.Add(a.parent.makeAuthorizationHeader())
	httpRequestHeaders.Add(iHttp.HttpHeaderClientTimeZoneName, TimeZoneForTests)

	var httpRequestUrl *url.URL
	httpRequestUrl, err = url.Parse(requestPath)
	if err != nil {
		return nil, err
	}

	httpRequestUrlQuery := httpRequestUrl.Query()
	httpRequestUrlQuery.Add(iHttp.HttpQueryParameterFileID, string(FileId))
	httpRequestUrl.RawQuery = httpRequestUrlQuery.Encode()

	httpRequest := &http.Request{
		URL:    httpRequestUrl,
		Method: http.MethodGet,
		Header: httpRequestHeaders,
	}

	var httpResponse *http.Response
	httpResponse, err = http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.Errorf(ErrFHttpStatusCodeWrong, http.StatusOK, httpResponse.StatusCode)
	}

	var httpResponseBody []byte
	httpResponseBody, err = io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		derr := httpResponse.Body.Close()
		err = multierr.Combine(err, derr)
	}()

	err = json.Unmarshal(httpResponseBody, &data.Data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// readLastEvents читает недавние события файла (последние N штук).
func (a *API) readLastEvents() (data *FileEventsLastNData, err error) {
	data = &FileEventsLastNData{}

	serverAddress := net.JoinHostPort(
		a.BusinessLogicsHttpServerSettings.HttpServerHost,
		strconv.Itoa(int(a.BusinessLogicsHttpServerSettings.HttpServerPort)),
	)

	requestPath := ProtocolSchemaFull + path.Join(serverAddress, "file-events", "last-n")

	httpRequestHeaders := make(http.Header)
	httpRequestHeaders.Add(a.parent.makeAuthorizationHeader())
	httpRequestHeaders.Add(iHttp.HttpHeaderClientTimeZoneName, TimeZoneForTests)

	var httpRequestUrl *url.URL
	httpRequestUrl, err = url.Parse(requestPath)
	if err != nil {
		return nil, err
	}

	httpRequestUrlQuery := httpRequestUrl.Query()
	httpRequestUrlQuery.Add(iHttp.HttpQueryParameterFileID, string(FileId))
	httpRequestUrl.RawQuery = httpRequestUrlQuery.Encode()

	httpRequest := &http.Request{
		URL:    httpRequestUrl,
		Method: http.MethodGet,
		Header: httpRequestHeaders,
	}

	var httpResponse *http.Response
	httpResponse, err = http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.Errorf(ErrFHttpStatusCodeWrong, http.StatusOK, httpResponse.StatusCode)
	}

	var httpResponseBody []byte
	httpResponseBody, err = io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		derr := httpResponse.Body.Close()
		err = multierr.Combine(err, derr)
	}()

	err = json.Unmarshal(httpResponseBody, &data.Data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
