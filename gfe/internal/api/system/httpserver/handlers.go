package httpserver

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	iHttp "github.com/vault-thirteen/junk/gfe/internal/http"
	"github.com/vault-thirteen/junk/gfe/internal/message"
	"github.com/vault-thirteen/junk/gfe/internal/storage"
)

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
