package models

import (
	"database/sql"
	"errors"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"time"
)

const (
	ErrUnexpectedNull = "unexpected null"
)

type Notification struct {
	// Identifier of this notification.
	Id base2.Id `json:"id"`

	// Identifier of a recipient user.
	UserId base2.Id `json:"userId"`

	// Textual information of this notification.
	Text base2.Text `json:"text"`

	// Time of creation.
	TimeOfCreation time.Time `json:"toc"`

	// Is the notification read by the recipient ?
	IsRead base2.Flag `json:"isRead"`

	// Time of reading.
	TimeOfReading *time.Time `json:"tor"`
}

func NewNotification() (n *Notification) {
	return &Notification{}
}

func NewNotificationFromScannableSource(src base.IScannable) (n *Notification, err error) {
	n = NewNotification()

	err = src.Scan(
		&n.Id,
		&n.UserId,
		&n.Text,
		&n.TimeOfCreation,
		&n.IsRead,
		&n.TimeOfReading,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return n, nil
}

func NewNotificationArrayFromRows(rows base.IScannableSequence) (ns []Notification, err error) {
	ns = []Notification{}
	var n *Notification

	for rows.Next() {
		n, err = NewNotificationFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		ns = append(ns, *n)
	}

	return ns, nil
}

func ListNotificationIds(notifications []Notification) (ids []base2.Id) {
	ids = make([]base2.Id, 0, len(notifications))

	for _, n := range notifications {
		ids = append(ids, n.Id)
	}

	return ids
}
