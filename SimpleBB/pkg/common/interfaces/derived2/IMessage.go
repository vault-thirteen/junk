package derived2

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"time"
)

type IMessage interface {
	GetLastTouchTime() time.Time

	// Emulated class members.
	GetIdPtr() (id *cmb.Id)
	GetThreadIdPtr() (threadId *cmb.Id)
	GetThreadId() (threadId cmb.Id)
	GetTextPtr() (text *cmb.Text)
	GetTextChecksumPtr() (textChecksum *[]byte)
	GetEventDataPtr() base2.IEventData
	GetEventData() base2.IEventData
	SetEventData(ed base2.IEventData)
	SetText(text cmb.Text)
}
