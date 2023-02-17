package random

import (
	"fmt"

	"github.com/google/uuid"
)

func MakeUniqueRandomString() string {
	return fmt.Sprintf("%s-%s", uuid.NewString(), uuid.NewString())
}
