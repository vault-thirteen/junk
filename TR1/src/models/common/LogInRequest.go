package cm

import "net"

type LogInRequest struct {
	MetaData
	Id int `gorm:"primarykey"`

	// Fields requested by a user.
	UserEmail string `gorm:"uniqueIndex,size:255"`

	// System fields.
	UserId           int    `gorm:"uniqueIndex"`
	RequestId        string `gorm:"uniqueIndex,size:255"`
	UserIPAB         net.IP
	CaptchaId        string
	VerificationCode string
	AuthData         []byte
}
