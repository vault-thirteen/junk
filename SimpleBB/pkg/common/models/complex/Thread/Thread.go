package t

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

type thread struct {
	// Identifier of this thread.
	Id cmb.Id `json:"id"`

	// Identifier of a forum containing this thread.
	ForumId cmb.Id `json:"forumId"`

	// Name of this thread.
	Name cms.Name `json:"name"`

	// List of identifiers of messages of this thread.
	Messages *ul.UidList `json:"messages"`

	// Thread meta-data.
	base2.IEventData
}

func NewThread() (t derived2.IThread) {
	return &thread{
		IEventData: ed.NewEventData(),
	}
}

func NewThreadFromScannableSource(src base.IScannable) (t derived2.IThread, err error) {
	t = NewThread()
	var x = ul.New()
	var eventData = ed.NewEventData()

	err = src.Scan(
		t.GetIdPtr(),
		t.GetForumIdPtr(),
		t.GetNamePtr(),
		x, //&t.Messages,
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

	t.SetEventData(eventData)
	t.SetMessages(x)
	return t, nil
}

func NewThreadArrayFromRows(rows base.IScannableSequence) (ts []derived2.IThread, err error) {
	ts = []derived2.IThread{}
	var t derived2.IThread

	for rows.Next() {
		t, err = NewThreadFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		ts = append(ts, t)
	}

	return ts, nil
}

// Emulated class members.
func (t *thread) GetIdPtr() (id *cmb.Id)                  { return &t.Id }
func (t *thread) GetForumIdPtr() (forumId *cmb.Id)        { return &t.ForumId }
func (t *thread) GetForumId() (forumId cmb.Id)            { return t.ForumId }
func (t *thread) GetNamePtr() (name *cms.Name)            { return &t.Name }
func (t *thread) GetMessagesPtr() (messages **ul.UidList) { return &t.Messages }
func (t *thread) GetMessages() (messages *ul.UidList)     { return t.Messages }
func (t *thread) GetEventDataPtr() base2.IEventData       { return t.IEventData }
func (t *thread) SetEventData(ed base2.IEventData) {
	t.IEventData = ed
}
func (t *thread) SetMessages(messages *ul.UidList) { t.Messages = messages }
