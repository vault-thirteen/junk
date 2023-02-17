package file

// ID -- идентификатор файла.
type ID string

// NewID -- конструктор ID.
// Конструктор возвращает не указатель на объект специально.
func NewID(s string) (id ID) {
	return ID(s)
}
