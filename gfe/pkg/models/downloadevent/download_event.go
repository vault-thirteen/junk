package downloadevent

import (
	"time"

	"github.com/vault-thirteen/gfe/pkg/models/user"
	"github.com/vault-thirteen/junk/gfe/pkg/models/file"
)

// DownloadEvent -- запись в базе данных о событии типа 'скачивание'.
// Используется в основном для тестов.
type DownloadEvent struct {
	// Идентификатор простого события.
	ID int `db:"id"`

	// Идентификатор пользователя.
	UserID user.ID `db:"user_id"`

	// Идентификатор файла.
	FileID file.ID `db:"file_id"`

	// Момент времени, в который произошло событие.
	// База данных хранит время событий в часовом поясе UTC.
	EventTime time.Time `db:"event_time"`
}
