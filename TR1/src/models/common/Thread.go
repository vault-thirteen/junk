package cm

type Thread struct {
	MetaData
	Id       int       `json:"id" gorm:"primarykey"`
	Name     string    `json:"name" gorm:"uniqueIndex,size:255"`
	ForumId  int       `json:"forumId"`
	Messages []Message `json:"messages,omitempty,omitzero"`
}
