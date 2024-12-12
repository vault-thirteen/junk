package derived2

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	ul "github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

type ISection interface {
	// Emulated class members.
	GetIdPtr() (id *cmb.Id)
	GetId() (id cmb.Id)
	GetParentPtr() (parent **cmb.Id)
	GetParent() (parent *cmb.Id)
	GetChildTypePtr() (childType derived1.ISectionChildType)
	GetChildType() (childType derived1.ISectionChildType)
	GetChildrenPtr() (children **ul.UidList)
	GetChildren() (children *ul.UidList)
	GetNamePtr() (name *cms.Name)
	GetEventDataPtr() base2.IEventData
	SetEventData(base2.IEventData)
	SetChildren(children *ul.UidList)
}
