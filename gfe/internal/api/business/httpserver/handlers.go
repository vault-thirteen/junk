package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kr/pretty"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	iHttp "github.com/vault-thirteen/junk/gfe/internal/http"
	"github.com/vault-thirteen/junk/gfe/internal/message"
	"github.com/vault-thirteen/junk/gfe/internal/storage"
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/history"
)

// MsgFDebugRequest -- сообщение для отладки запроса.
const MsgFDebugRequest = "request: %v"

// handlerLiveness -- обработчик, отвечающий о том, что сервис жив (существует
// и принимает запросы).
func (hs *HttpServer) handlerLiveness(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
}

// handlerReadiness -- обработчик, отвечающий о том, что сервис функционирует
// нормально.
func (hs *HttpServer) handlerReadiness(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	kafkaIsReady := hs.kafka.GetReadinessState()
	storageIsReady := hs.storage.IsReady()

	if kafkaIsReady && storageIsReady {
		return
	}

	// Проверяем доступность читателя Kafka.
	if !kafkaIsReady {
		iHttp.RespondWithNotReadyStatus(hs.logger, w, message.MsgKafkaNotReady)
		return
	}

	// Проверяем доступность хранилища.
	// Сначала -- методом грубой оценки, затем -- пингуем базу данных.
	if !storageIsReady {
		iHttp.RespondWithNotReadyStatus(hs.logger, w, storage.MsgStorageNotReady)
		return
	}

	var err = hs.storage.Ping()
	if err != nil {
		hs.logger.Err(err).Send()
		iHttp.RespondWithNotReadyStatus(hs.logger, w, storage.MsgStorageNotReady)
		return
	}

	iHttp.RespondWithNotReadyStatus(hs.logger, w, message.MsgSomethingNotReady)
}

// handlerGetAllEventsForFile -- обработчик, отдающий всю историю событий по
// файлу.
func (hs *HttpServer) handlerGetAllEventsForFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := hs.getRequestToGetFileEvents(r)
	if err != nil {
		iHttp.RespondWithBadRequestError(w, err.Error())
		return
	}

	req.RecordsCountLimit = 0

	hs.logger.Debug().Msg(pretty.Sprintf(MsgFDebugRequest, req))

	ctx, cancelFunc := context.WithTimeout(context.Background(), config.StorageQueryTimeoutSec*time.Second)
	defer cancelFunc()

	var events []*event.Event
	events, err = hs.storage.GetFileEvents(ctx, req)
	if err != nil {
		iHttp.RespondWithInternalServerError(w, err.Error())
		return
	}

	var response = &history.History{
		FileID:       req.FileID,
		TimeZoneName: req.ClientTimeZone,
		Records:      events,
	}

	iHttp.RespondWithJsonObject(hs.logger, w, response)
}

// handlerGetLastNEventsForFile -- обработчик, отдающий N последних записей из
// истории событий по файлу.
func (hs *HttpServer) handlerGetLastNEventsForFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := hs.getRequestToGetFileEvents(r)
	if err != nil {
		iHttp.RespondWithBadRequestError(w, err.Error())
		return
	}

	req.RecordsCountLimit = config.FileLastEventsCount

	hs.logger.Debug().Msg(pretty.Sprintf(MsgFDebugRequest, req))

	ctx, cancelFunc := context.WithTimeout(context.Background(), config.StorageQueryTimeoutSec*time.Second)
	defer cancelFunc()

	var events []*event.Event
	events, err = hs.storage.GetFileEvents(ctx, req)
	if err != nil {
		iHttp.RespondWithInternalServerError(w, err.Error())
		return
	}

	var response = &history.History{
		FileID:       req.FileID,
		TimeZoneName: req.ClientTimeZone,
		Records:      events,
	}

	iHttp.RespondWithJsonObject(hs.logger, w, response)
}

// handlerGetFileEventTypes -- обработчик, отдающий список типов событий.
func (hs *HttpServer) handlerGetFileEventTypes(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), config.StorageQueryTimeoutSec*time.Second)
	defer cancelFunc()

	// Для простоты, читаем список из хранилища.
	// Если запросов будет много, то можно будет приделать кэш.
	eventTypes, err := hs.storage.GetFileEventTypes(ctx)
	if err != nil {
		iHttp.RespondWithInternalServerError(w, err.Error())
		return
	}

	iHttp.RespondWithJsonObject(hs.logger, w, eventTypes)
}
