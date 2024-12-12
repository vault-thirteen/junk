package derived2

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

type IForum interface {
	// Emulated class members.
	GetIdPtr() (id *cmb.Id)
	GetId() (id cmb.Id)
	GetSectionIdPtr() (sectionId *cmb.Id)
	GetSectionId() (sectionId cmb.Id)
	GetNamePtr() (name *cms.Name)
	GetThreadsPtr() (threads **ul.UidList)
	GetThreads() (threads *ul.UidList)
	GetEventDataPtr() base2.IEventData
	SetEventData(ed base2.IEventData)
	SetThreads(threads *ul.UidList)
}
