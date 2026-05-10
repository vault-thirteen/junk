package cm

import "net"

type LogOutRequest struct {
	MetaData
	Id int `gorm:"primarykey"`

	// System fields.
	UserId    int    `gorm:"uniqueIndex"`
	RequestId string `gorm:"uniqueIndex,size:255"`
	UserIPAB  net.IP
}
