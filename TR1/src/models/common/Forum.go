package cm

type Forum struct {
	MetaData
	Id      int      `json:"id" gorm:"primarykey"`
	Name    string   `json:"name" gorm:"uniqueIndex,size:255"`
	Threads []Thread `json:"threads,omitempty,omitzero"`

	// Forum's position in the list of all forums.
	// 1 = top position; the greater number is, the lower the position is.
	Pos int `json:"pos" gorm:"uniqueIndex"`
}
