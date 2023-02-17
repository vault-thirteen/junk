package random

import (
	"fmt"

	"github.com/google/uuid"
)

// MakeUniqueRandomString создаёт уникальную случайную строку.
func MakeUniqueRandomString() string {
	return fmt.Sprintf("%s-%s", uuid.NewString(), uuid.NewString())
}

//TODO:
//	- Получить информацию по Vault.
//	- Проверить работу Kafka Consumer Group при нескольких запущенных экземплярах приложения.
//	- Нужен ли уникальный ID инстантса сервиса для чего-либо ?
