package base

import (
	"database/sql/driver"
	sql2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
	"strconv"
)

type Count int

func (c *Count) Scan(src any) (err error) {
	var x int
	x, err = sql2.ScanSrcAsInt(src)
	if err != nil {
		return err
	}

	*c = Count(x)
	return nil
}

func (c Count) Value() (dv driver.Value, err error) {
	return driver.Value(sql2.IntToSql(c.AsInt())), nil
}

func (c Count) ToString() string {
	return strconv.Itoa(c.AsInt())
}

func (c Count) RawValue() int {
	return c.AsInt()
}

func (c Count) AsInt() int {
	return int(c)
}
