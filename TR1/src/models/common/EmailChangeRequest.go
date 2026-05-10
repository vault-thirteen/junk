package cm

import "net"

type EmailChangeRequest struct {
	MetaData
	Id int `gorm:"primarykey"`

	// Fields requested by a user.
	NewEmail string `gorm:"uniqueIndex,size:255"`

	// System fields.
	UserId            int    `gorm:"uniqueIndex"`
	RequestId         string `gorm:"uniqueIndex,size:255"`
	UserIPAB          net.IP
	CaptchaId         string
	VerificationCodeA string
	VerificationCodeB string
	AuthData          []byte
}
