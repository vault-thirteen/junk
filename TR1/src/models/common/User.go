package cm

import (
	"net/mail"
	"time"

	"github.com/vault-thirteen/BytePackedPassword"
)

type User struct {
	MetaData
	Id                int        `json:"id" gorm:"primarykey"`
	Name              string     `json:"name,omitempty,omitzero" gorm:"uniqueIndex,size:255"`
	Email             string     `json:"email,omitempty,omitzero" gorm:"uniqueIndex,size:255"`
	Password          *Password  `json:"password,omitempty,omitzero"`
	Session           *Session   `json:"session,omitempty,omitzero"`
	Roles             *Roles     `json:"roles,omitempty,omitzero" gorm:"embedded"`
	RegTime           time.Time  `json:"regTime,omitempty,omitzero"`
	LastBadActionTime *time.Time `json:"lastBadActionTime,omitempty,omitzero"`
	BanTime           *time.Time `json:"banTime,omitempty,omitzero"`
}

func IsUserEmailValid(email string) (isValid bool) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}

func IsUserPasswordAllowed(password string) (isAllowed bool) {
	ok, _ := bpp.IsPasswordAllowed(password)
	return ok
}
