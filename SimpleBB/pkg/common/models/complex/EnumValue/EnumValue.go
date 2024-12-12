package ev

import (
	"database/sql/driver"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	json2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/json"
	sql2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
	"strconv"
)

// EnumValue is a raw type of the enumeration. It does not restrict the value.
// Validation of the value is performed by the 'Enum' class according to it's
// settings.
type enumValue byte

func NewEnumValue(b byte) cmi.IEnumValue {
	x := enumValue(b)
	return &x
}

func (ev *enumValue) Scan(src any) (err error) {
	var b byte
	b, err = sql2.ScanSrcAsByte(src)
	if err != nil {
		return err
	}

	*ev = enumValue(b)
	return nil
}

func (ev enumValue) Value() (dv driver.Value, err error) {
	return driver.Value(sql2.ByteToSql(ev.AsByte())), nil
}

func (ev enumValue) ToString() string {
	return strconv.Itoa(ev.AsInt())
}

func (ev enumValue) RawValue() byte {
	return ev.AsByte()
}

func (ev enumValue) AsByte() byte {
	return byte(ev)
}

func (ev enumValue) AsInt() int {
	return int(ev)
}

func (ev *enumValue) UnmarshalJSON(data []byte) (err error) {
	var b byte
	b, err = json2.UnmarshalAsByte(data)
	if err != nil {
		return err
	}

	*ev = enumValue(b)
	return nil
}

func (ev enumValue) MarshalJSON() (ba []byte, err error) {
	return json2.ToJson(ev), nil
}
