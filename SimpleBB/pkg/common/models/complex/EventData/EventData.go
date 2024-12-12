package ed

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"time"
)

type eventData struct {
	// Parameters of creation.
	Creator *simple.EventParameters `json:"creator"`

	// Parameters of the last edit.
	Editor *simple.OptionalEventParameters `json:"editor"`
}

func NewEventData() base2.IEventData {
	return &eventData{
		Creator: &simple.EventParameters{},
		Editor:  &simple.OptionalEventParameters{},
	}
}

// Emulated class members.
func (ed *eventData) GetCreatorUserIdPtr() (userId *cmb.Id) { return &ed.Creator.UserId }
func (ed *eventData) GetCreatorUserId() (userId cmb.Id)     { return ed.Creator.UserId }
func (ed *eventData) GetCreatorTimePtr() (time *time.Time)  { return &ed.Creator.Time }
func (ed *eventData) GetCreatorTime() (time time.Time)      { return ed.Creator.Time }
func (ed *eventData) GetEditorUserIdPtr() (userId **cmb.Id) { return &ed.Editor.UserId }
func (ed *eventData) GetEditorTimePtr() (time **time.Time)  { return &ed.Editor.Time }
func (ed *eventData) GetEditorTime() (time *time.Time)      { return ed.Editor.Time }
