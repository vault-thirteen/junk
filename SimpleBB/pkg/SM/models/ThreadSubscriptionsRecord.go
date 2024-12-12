package models

import (
	"database/sql"
	"errors"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type ThreadSubscriptionsRecord struct {
	Id       cmb.Id      `json:"id"`
	ThreadId cmb.Id      `json:"threadId"`
	Users    *ul.UidList `json:"userIds"`
}

func NewThreadSubscriptionsRecord() (tsr *ThreadSubscriptionsRecord) {
	return &ThreadSubscriptionsRecord{}
}

func NewThreadSubscriptionsRecordFromScannableSource(src base.IScannable) (tsr *ThreadSubscriptionsRecord, err error) {
	tsr = NewThreadSubscriptionsRecord()
	var x = ul.New()

	err = src.Scan(
		&tsr.Id,
		&tsr.ThreadId,
		x, //&tsr.Users,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	tsr.Users = x
	return tsr, nil
}

func NewThreadSubscriptionsRecordArrayFromRows(rows base.IScannableSequence) (tsrs []ThreadSubscriptionsRecord, err error) {
	tsrs = []ThreadSubscriptionsRecord{}
	var tsr *ThreadSubscriptionsRecord

	for rows.Next() {
		tsr, err = NewThreadSubscriptionsRecordFromScannableSource(rows)
		if err != nil {
			return nil, err
		}

		tsrs = append(tsrs, *tsr)
	}

	return tsrs, nil
}
