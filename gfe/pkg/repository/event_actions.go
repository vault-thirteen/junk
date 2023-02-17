package repository

import (
	"context"

	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	fer "github.com/vault-thirteen/junk/gfe/pkg/models/fileeventsrequest"
	"github.com/vault-thirteen/junk/gfe/pkg/models/message"
)

// EventActions описывает контракт взаимодействия
// с хранилищем событий по файлам.
type EventActions interface {
	// GetFileEventTypes читает из базы данных типы файловых событий.
	GetFileEventTypes(ctx context.Context) (eventTypes []event.Type, err error)

	// SaveEvent сохраняет в базу данных файловое событие.
	SaveEvent(ctx context.Context, eventMessage *message.Message) (err error)

	// GetFileEvents читает из базы данных файловые события.
	GetFileEvents(ctx context.Context, req *fer.FileEventsRequest) (events []*event.Event, err error)
}
