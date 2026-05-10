package cm

import (
	"net"
)

type Session struct {
	MetaData
	Id       int    `json:"id" gorm:"primarykey"`
	UserId   int    `json:"userId" gorm:"uniqueIndex"`
	User     *User  `json:"user"`
	UserIPAB net.IP `json:"userIPAB"`
}
