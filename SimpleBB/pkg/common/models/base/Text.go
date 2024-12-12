package base

import (
	"database/sql/driver"
	sql2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
)

type Text string

func (t *Text) Scan(src any) (err error) {
	var s string
	s, err = sql2.ScanSrcAsString(src)
	if err != nil {
		return err
	}

	*t = Text(s)
	return nil
}

func (t Text) Value() (dv driver.Value, err error) {
	return driver.Value(sql2.StringToSql(t.AsString())), nil
}

func (t Text) ToString() string {
	return t.AsString()
}

func (t Text) RawValue() string {
	return t.AsString()
}

func (t Text) AsString() string {
	return string(t)
}

func (t Text) AsBytes() []byte {
	return []byte(t)
}
