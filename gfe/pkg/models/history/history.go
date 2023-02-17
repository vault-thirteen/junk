package history

import (
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/file"
)

// History -- история событий по файлу.
// Используется для HTTP запросов.
type History struct {
	// Идентификатор файла.
	FileID file.ID `json:"fileId"`

	// Запрошенный часовой пояс.
	TimeZoneName string `json:"timeZone"`

	// Список событий по запрошенному файлу с учётом запрошенного часового
	// пояса.
	Records []*event.Event `json:"records"`
}
