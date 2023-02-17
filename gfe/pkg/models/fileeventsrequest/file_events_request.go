package fileeventsrequest

import "github.com/vault-thirteen/junk/gfe/pkg/models/file"

// FileEventsRequest хранит параметры запроса событий по файлу.
type FileEventsRequest struct {
	// Идентификатор файла.
	FileID file.ID

	// Название часового пояса клиента.
	ClientTimeZone string

	// Ограничитель количества записей в ответе.
	// Если значение менее 1, то ограничитель не применяется.
	RecordsCountLimit int
}

// IsRecordsCountLimitSet возвращает true если запрос предполагает ограничение
// количества записей; возвращает false в противном случае.
func (r *FileEventsRequest) IsRecordsCountLimitSet() bool {
	if r.RecordsCountLimit < 1 {
		return false
	}

	return true
}
