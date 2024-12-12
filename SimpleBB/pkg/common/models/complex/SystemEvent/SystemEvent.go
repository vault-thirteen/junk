package se

import (
	"database/sql"
	"errors"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SystemEventData"
	"time"
)

type systemEvent struct {
	derived2.ISystemEventData

	// ID and time of the event are automatically set by database and should
	// not be touched manually.
	Id   cmb.Id    `json:"id"`
	Time time.Time `json:"time"`
}

func NewSystemEvent() (se derived2.ISystemEvent) {
	return &systemEvent{
		ISystemEventData: sed.NewSystemEventData(),
	}
}

func NewSystemEventWithData(data derived2.ISystemEventData) (se derived2.ISystemEvent, err error) {
	se = &systemEvent{
		ISystemEventData: data,
	}

	_, err = se.GetSystemEventData().CheckParameters()
	if err != nil {
		return nil, err
	}

	return se, nil
}

func NewSystemEventFromScannableSource(src base.IScannable) (se derived2.ISystemEvent, err error) {
	se = NewSystemEvent()

	x := se.GetSystemEventData()

	err = src.Scan(
		se.GetIdPtr(),
		x.GetThreadIdPtr(),
		x.GetMessageIdPtr(),
		x.GetUserIdPtr(),
		se.GetTimePtr(),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return se, nil
}

func NewSystemEventArrayFromRows(rows base.IScannableSequence) (ses []derived2.ISystemEvent, err error) {
	ses = []derived2.ISystemEvent{}
	var se derived2.ISystemEvent

	for rows.Next() {
		se, err = NewSystemEventFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		ses = append(ses, se)
	}

	return ses, nil
}

// Emulated class members.
func (se *systemEvent) GetIdPtr() (id *cmb.Id) {
	return &se.Id
}
func (se *systemEvent) GetTimePtr() (t *time.Time) { return &se.Time }
func (se *systemEvent) GetSystemEventData() (sed derived2.ISystemEventData) {
	return se.ISystemEventData
}
