package cm

import (
	"net"
)

const (
	LogEvent_Type_LogIn                 = 1
	LogEvent_Type_LogOutBySelf          = 2
	LogEvent_Type_LogOutByTimeout       = 3
	LogEvent_Type_LogOutByAction        = 4
	LogEvent_Type_LogOutByAdministrator = 5
)

type LogEvent struct {
	MetaData
	Id       int `gorm:"primarykey"`
	UserId   int `gorm:"index"`
	UserIPAB net.IP
	Type     int
	AdminId  *int
}

func NewLogEvent(typé int, userId int, userIPAB net.IP, adminId *int) *LogEvent {
	return &LogEvent{
		UserId:   userId,
		UserIPAB: userIPAB,
		Type:     typé,
		AdminId:  adminId,
	}
}
