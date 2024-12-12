package derived2

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"time"
)

type IResource interface {
	IsValid() bool
	AsText() cmb.Text
	AsNumber() cmb.Count
	GetValue() any

	// Emulated class members.
	GetIdPtr() (id *cmb.Id)
	GetTypePtr() (t derived1.IResourceType)
	GetTextPtr() (text **cmb.Text)
	GetNumberPtr() (number **cmb.Count)
	GetTocPtr() (toc *time.Time)
	GetType() (t derived1.IResourceType)
	GetText() (text *cmb.Text)
	GetNumber() (number *cmb.Count)
}
