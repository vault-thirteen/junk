package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	fer "github.com/vault-thirteen/junk/gfe/pkg/models/fileeventsrequest"
	"github.com/vault-thirteen/junk/gfe/pkg/models/message"
)

// Сообщения.
const (
	msgAggregatedEvent = "aggregated event"
	msgSimpleEvent     = "simple event"
)

// Ошибки.
const (
	// ErrFInsertedRowsCount -- сообщение об ошибке "количество вставленных рядов (строк) равно".
	ErrFInsertedRowsCount = "inserted rows count is %v"
)

// GetFileEventTypes делает из базы данных выборку типов событий.
func (s *Storage) GetFileEventTypes(ctx context.Context) (eventTypes []event.Type, err error) {
	s.sqlQueryPreparedStatementsLock.RLock()
	defer s.sqlQueryPreparedStatementsLock.RUnlock()

	st := s.sqlQueryPreparedStatements[SqlQueryIndexForQueryToGetFileEventTypes]

	err = st.SelectContext(ctx, &eventTypes)
	if err != nil {
		return nil, err
	}

	return eventTypes, nil
}

// SaveEvent сохраняет событие в хранилище (базу данных).
func (s *Storage) SaveEvent(ctx context.Context, eventMessage *message.Message) (err error) {
	if eventMessage == nil {
		return message.ErrEventMessageNull
	}

	_, err = eventMessage.EventTypeID.IsValid()
	if err != nil {
		return err
	}

	if eventMessage.EventTypeID.IsAggregated() {
		s.logger.Debug().Msg(msgAggregatedEvent)

		switch eventMessage.EventTypeID {
		case event.Download:
			err = s.saveDownloadEvent(ctx, eventMessage)
			if err != nil {
				return err
			}

			return nil
		}
	}

	if eventMessage.EventTypeID.IsSimple() {
		s.logger.Debug().Msg(msgSimpleEvent)

		switch eventMessage.EventTypeID {
		case event.Creation:
			err = s.saveCreationEvent(ctx, eventMessage)
			if err != nil {
				return err
			}

			return nil

		case event.Upload:
			err = s.saveUploadEvent(ctx, eventMessage)
			if err != nil {
				return err
			}

			return nil

		case event.Modification:
			err = s.saveModificationEvent(ctx, eventMessage)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.Errorf(event.ErrFTypeUnsupported, eventMessage.EventTypeID)
}

// saveCreationEvent сохраняет в хранилище событие типа "создание" (ID=1).
func (s *Storage) saveCreationEvent(ctx context.Context, eventMessage *message.Message) (err error) {
	return s.saveSimpleEvent(ctx, eventMessage)
}

// saveUploadEvent сохраняет в хранилище событие типа "загрузка" (ID=2).
func (s *Storage) saveUploadEvent(ctx context.Context, eventMessage *message.Message) (err error) {
	return s.saveSimpleEvent(ctx, eventMessage)
}

// saveDownloadEvent сохраняет в хранилище событие типа "скачивание" (ID=3).
func (s *Storage) saveDownloadEvent(ctx context.Context, eventMessage *message.Message) (err error) {
	s.sqlQueryPreparedStatementsLock.RLock()
	defer s.sqlQueryPreparedStatementsLock.RUnlock()

	st := s.sqlQueryPreparedStatements[SqlQueryIndexForQueryToInsertFileDownloadEvent]

	var sqlResult sql.Result
	sqlResult, err = st.ExecContext(
		ctx,
		eventMessage.UserID,
		eventMessage.FileID,
		eventMessage.EventTime)
	if err != nil {
		return err
	}

	err = s.checkInsertion(sqlResult)
	if err != nil {
		return err
	}

	return nil
}

// saveModificationEvent сохраняет в хранилище событие типа "изменение" (ID=4).
func (s *Storage) saveModificationEvent(ctx context.Context, eventMessage *message.Message) (err error) {
	return s.saveSimpleEvent(ctx, eventMessage)
}

// saveSimpleEvent сохраняет в хранилище простое событие (ID in {1,2,4}).
// Параметр "eventMessage.ID" должен быть правильным!
func (s *Storage) saveSimpleEvent(ctx context.Context, eventMessage *message.Message) (err error) {
	s.sqlQueryPreparedStatementsLock.RLock()
	defer s.sqlQueryPreparedStatementsLock.RUnlock()

	st := s.sqlQueryPreparedStatements[SqlQueryIndexForQueryToInsertFileSimpleEvent]

	var sqlResult sql.Result
	sqlResult, err = st.ExecContext(
		ctx,
		eventMessage.UserID,
		eventMessage.FileID,
		eventMessage.EventTypeID,
		eventMessage.EventTime,
	)
	if err != nil {
		return err
	}

	err = s.checkInsertion(sqlResult)
	if err != nil {
		return err
	}

	return nil
}

// checkInsertion проверяет, что в таблицу была вставлена ровно одна строка.
func (s *Storage) checkInsertion(sqlResult sql.Result) (err error) {
	var insertedRowsCount int64
	insertedRowsCount, err = sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if insertedRowsCount != 1 {
		return errors.Errorf(ErrFInsertedRowsCount, insertedRowsCount)
	}

	return nil
}

// GetFileEvents делает из базы данных выборку событий по файлу.
func (s *Storage) GetFileEvents(ctx context.Context, req *fer.FileEventsRequest) (events []*event.Event, err error) {
	s.sqlQueryPreparedStatementsLock.RLock()
	defer s.sqlQueryPreparedStatementsLock.RUnlock()

	var st *sqlx.Stmt
	if req.IsRecordsCountLimitSet() {
		st = s.sqlQueryPreparedStatements[SqlQueryIndexForQueryToSelectFileEventsLimited]
	} else {
		st = s.sqlQueryPreparedStatements[SqlQueryIndexForQueryToSelectFileEventsAll]
	}

	queryParameters := s.getArgumentsForQuerySelectFileEvents(req)

	var rawRecords = make([]*event.RawEvent, 0)

	err = st.SelectContext(ctx, &rawRecords, queryParameters...)
	if err != nil {
		return nil, err
	}

	events = make([]*event.Event, 0, len(rawRecords))
	var ev *event.Event
	for _, rawRecord := range rawRecords {
		ev, err = event.NewFromRawData(rawRecord, req.ClientTimeZone)
		if err != nil {
			return nil, err
		}

		events = append(events, ev)
	}

	return events, nil
}

// getArgumentsForQuerySelectFileEvents возвращает набор аргументов для запроса
// на выборку из базы данных событий по файлу.
func (s *Storage) getArgumentsForQuerySelectFileEvents(req *fer.FileEventsRequest) (queryArguments []interface{}) {
	// Если отдаём все записи, то не прикрепляем секцию 'LIMIT'.
	if !req.IsRecordsCountLimitSet() {
		return []interface{}{
			req.ClientTimeZone, // $1.
			req.ClientTimeZone, // $2.
			req.FileID,         // $3.
			req.ClientTimeZone, // $4.
			req.FileID,         // $5.
		}
	}

	// Используем ограничитель.
	return []interface{}{
		req.ClientTimeZone,    // $1.
		req.ClientTimeZone,    // $2.
		req.FileID,            // $3.
		req.ClientTimeZone,    // $4.
		req.FileID,            // $5.
		req.RecordsCountLimit, // $6.
	}
}
