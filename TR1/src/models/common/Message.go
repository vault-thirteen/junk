package cm

type Message struct {
	MetaData
	Id       int    `json:"id" gorm:"primarykey"`
	Text     string `json:"text"`
	ThreadId int    `json:"threadId"`

	// This behaviour is not clearly documented !
	// Source: https://github.com/go-gorm/gorm/issues/6463
	Creator   *User `json:"creator,omitempty,omitzero" gorm:"foreignKey:CreatorId"`
	CreatorId int   `json:"creatorId"`

	// This behaviour is not clearly documented !
	// Source: https://github.com/go-gorm/gorm/issues/6463
	Editor   *User `json:"editor,omitempty,omitzero" gorm:"foreignKey:EditorId"`
	EditorId *int  `json:"editorId"`
}
