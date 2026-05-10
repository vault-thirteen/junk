package cm

type Password struct {
	MetaData
	Id     int `gorm:"primarykey"`
	UserId int `gorm:"uniqueIndex"`
	Bytes  []byte
	Text   string
}
