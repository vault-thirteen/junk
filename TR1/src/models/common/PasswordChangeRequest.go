package cm

import "net"

type PasswordChangeRequest struct {
	MetaData
	Id int `gorm:"primarykey"`

	// Fields requested by a user.
	NewPassword string `gorm:"uniqueIndex,size:255"`

	// System fields.
	UserId           int    `gorm:"uniqueIndex"`
	RequestId        string `gorm:"uniqueIndex,size:255"`
	UserIPAB         net.IP
	CaptchaId        string
	VerificationCode string
	AuthData         []byte
}
