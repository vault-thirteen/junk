package derived2

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

type IThread interface {
	// Emulated class members.
	GetIdPtr() (id *cmb.Id)
	GetForumIdPtr() (forumId *cmb.Id)
	GetForumId() (forumId cmb.Id)
	GetNamePtr() (name *cms.Name)
	GetMessagesPtr() (messages **ul.UidList)
	GetMessages() (messages *ul.UidList)
	GetEventDataPtr() base2.IEventData
	SetEventData(ed base2.IEventData)
	SetMessages(messages *ul.UidList)
}
