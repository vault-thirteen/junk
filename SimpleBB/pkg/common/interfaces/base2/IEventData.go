package base2

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"time"
)

type IEventData interface {
	// Emulated class members.
	GetCreatorUserIdPtr() (userId *cmb.Id)
	GetCreatorUserId() (userId cmb.Id)
	GetCreatorTimePtr() (time *time.Time)
	GetCreatorTime() (time time.Time)
	GetEditorUserIdPtr() (userId **cmb.Id)
	GetEditorTimePtr() (time **time.Time)
	GetEditorTime() (time *time.Time)
}
