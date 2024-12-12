package f

import (
	"database/sql"
	"errors"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	ed "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EventData"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

type forum struct {
	// Identifier of this forum.
	Id cmb.Id `json:"id"`

	// Identifier of a section containing this forum.
	SectionId cmb.Id `json:"sectionId"`

	// Name of this forum.
	Name cms.Name `json:"name"`

	// List of identifiers of threads of this forum.
	Threads *ul.UidList `json:"threads"`

	// Forum meta-data.
	base2.IEventData
}

func NewForum() (frm derived2.IForum) {
	return &forum{
		IEventData: ed.NewEventData(),
	}
}

func NewForumFromScannableSource(src base.IScannable) (forum derived2.IForum, err error) {
	forum = NewForum()
	var x = ul.New()
	var eventData = ed.NewEventData()

	err = src.Scan(
		forum.GetIdPtr(),
		forum.GetSectionIdPtr(),
		forum.GetNamePtr(),
		x, //&forum.Threads,
		eventData.GetCreatorUserIdPtr(),
		eventData.GetCreatorTimePtr(),
		eventData.GetEditorUserIdPtr(),
		eventData.GetEditorTimePtr(),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	forum.SetEventData(eventData)
	forum.SetThreads(x)
	return forum, nil
}

func NewForumArrayFromRows(rows base.IScannableSequence) (forums []derived2.IForum, err error) {
	forums = []derived2.IForum{}
	var f derived2.IForum

	for rows.Next() {
		f, err = NewForumFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		forums = append(forums, f)
	}

	return forums, nil
}

// Emulated class members.
func (f *forum) GetIdPtr() (id *cmb.Id)                { return &f.Id }
func (f *forum) GetId() (id cmb.Id)                    { return f.Id }
func (f *forum) GetSectionIdPtr() (sectionId *cmb.Id)  { return &f.SectionId }
func (f *forum) GetSectionId() (sectionId cmb.Id)      { return f.SectionId }
func (f *forum) GetNamePtr() (name *cms.Name)          { return &f.Name }
func (f *forum) GetThreadsPtr() (threads **ul.UidList) { return &f.Threads }
func (f *forum) GetThreads() (threads *ul.UidList)     { return f.Threads }
func (f *forum) GetEventDataPtr() base2.IEventData     { return f.IEventData }
func (f *forum) SetEventData(ed base2.IEventData) {
	f.IEventData = ed
}
func (f *forum) SetThreads(threads *ul.UidList) { f.Threads = threads }
