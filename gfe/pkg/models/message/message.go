package message

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/file"
	"github.com/vault-thirteen/junk/gfe/pkg/models/user"
)

var (
	// ErrEventMessageNull -- сообщение об ошибке "сообщение о событии не установлено".
	ErrEventMessageNull = errors.New("event message is not set")
)

// Message -- сообщение о произошедшем событии.
// Используется для Kafka.
type Message struct {
	// Идентификатор пользователя.
	UserID user.ID `json:"userId"`

	// Идентификатор файла.
	FileID file.ID `json:"fileId"`

	// Идентификатор типа события.
	EventTypeID event.TypeID `json:"eventTypeId"`

	// Момент времени, в который произошло событие.
	// База данных хранит время событий в часовом поясе UTC.
	EventTime time.Time `json:"eventTime"`

	// Исходное сообщение от Kafka.
	RawMessage *sarama.ConsumerMessage `json:"-"`
}
