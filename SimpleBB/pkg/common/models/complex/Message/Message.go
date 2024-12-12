package m

import (
	"database/sql"
	"errors"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	ed "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EventData"
	"time"
)

type message struct {
	// Identifier of this message.
	Id cmb.Id `json:"id"`

	// Identifier of a thread containing this message.
	ThreadId cmb.Id `json:"threadId"`

	// Textual information of this message.
	Text cmb.Text `json:"text"`

	// Check sum of the Text field.
	TextChecksum []byte `json:"textChecksum"`

	// Message meta-data.
	base2.IEventData
}

func NewMessage() (msg derived2.IMessage) {
	return &message{
		IEventData: ed.NewEventData(),
	}
}

func NewMessageFromScannableSource(src base.IScannable) (msg derived2.IMessage, err error) {
	msg = NewMessage()
	var eventData = ed.NewEventData()

	err = src.Scan(
		msg.GetIdPtr(),
		msg.GetThreadIdPtr(),
		msg.GetTextPtr(),
		msg.GetTextChecksumPtr(),
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

	msg.SetEventData(eventData)
	return msg, nil
}

func NewMessageArrayFromRows(rows base.IScannableSequence) (msgs []derived2.IMessage, err error) {
	msgs = []derived2.IMessage{}
	var msg derived2.IMessage

	for rows.Next() {
		msg, err = NewMessageFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (m *message) GetLastTouchTime() time.Time {
	if m.IEventData.GetEditorTime() == nil {
		return m.GetCreatorTime()
	}

	return *m.GetEditorTime()
}

// Emulated class members.
func (m *message) GetIdPtr() (id *cmb.Id)                     { return &m.Id }
func (m *message) GetThreadIdPtr() (threadId *cmb.Id)         { return &m.ThreadId }
func (m *message) GetThreadId() (threadId cmb.Id)             { return m.ThreadId }
func (m *message) GetTextPtr() (text *cmb.Text)               { return &m.Text }
func (m *message) GetTextChecksumPtr() (textChecksum *[]byte) { return &m.TextChecksum }
func (m *message) GetEventDataPtr() base2.IEventData          { return m.IEventData }
func (m *message) GetEventData() base2.IEventData             { return m.IEventData }
func (m *message) SetEventData(ed base2.IEventData) {
	m.IEventData = ed
}
func (m *message) SetText(text cmb.Text) { m.Text = text }
