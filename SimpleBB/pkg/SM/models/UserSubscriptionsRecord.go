package models

import (
	"database/sql"
	"errors"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	IdForVirtualUserSubscriptionsRecord = -1
)

type UserSubscriptionsRecord struct {
	Id      cmb.Id      `json:"id"`
	UserId  cmb.Id      `json:"userId"`
	Threads *ul.UidList `json:"threadIds"`
}

func NewUserSubscriptionsRecord() (usr *UserSubscriptionsRecord) {
	return &UserSubscriptionsRecord{}
}

func NewUserSubscriptionsRecordFromScannableSource(src base.IScannable) (usr *UserSubscriptionsRecord, err error) {
	usr = NewUserSubscriptionsRecord()
	var x = new(ul.UidList)

	err = src.Scan(
		&usr.Id,
		&usr.UserId,
		x, //&us.Threads,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	usr.Threads = x
	return usr, nil
}

func NewUserSubscriptionsRecordArrayFromRows(rows base.IScannableSequence) (usrs []UserSubscriptionsRecord, err error) {
	usrs = []UserSubscriptionsRecord{}
	var usr *UserSubscriptionsRecord

	for rows.Next() {
		usr, err = NewUserSubscriptionsRecordFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		usrs = append(usrs, *usr)
	}

	return usrs, nil
}
