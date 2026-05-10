package cm

import "time"

type MetaData struct {
	// Fields for GORM.
	CreatedAt time.Time `json:"toc,omitempty,omitzero"`
	UpdatedAt time.Time `json:"tou,omitempty,omitzero"`
}
