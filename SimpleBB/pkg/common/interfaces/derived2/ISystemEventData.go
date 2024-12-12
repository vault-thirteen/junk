package derived2

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type ISystemEventData interface {
	CheckParameters() (ok bool, err error)

	// Emulated class members.
	GetType() (t derived1.ISystemEventType)
	GetThreadId() (threadId *cmb.Id)
	GetThreadIdPtr() (threadId **cmb.Id)
	GetMessageIdPtr() (messageId **cmb.Id)
	GetMessageId() (messageId *cmb.Id)
	GetUserIdPtr() (userId **cmb.Id)
	GetUserId() (userId *cmb.Id)
	GetCreator() (creator *cmb.Id)
}
