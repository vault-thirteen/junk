package simpleevent

import (
	"time"

	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/file"
	"github.com/vault-thirteen/junk/gfe/pkg/models/user"
)

// SimpleEvent -- запись в базе данных о простом событии.
// Используется в основном для тестов.
type SimpleEvent struct {
	// Идентификатор простого события.
	ID int `db:"id"`

	// Идентификатор пользователя.
	UserID user.ID `db:"user_id"`

	// Идентификатор файла.
	FileID file.ID `db:"file_id"`

	// Идентификатор типа события.
	EventTypeID event.TypeID `db:"event_type_id"`

	// Момент времени, в который произошло событие.
	// База данных хранит время событий в часовом поясе UTC.
	EventTime time.Time `db:"event_time"`
}
