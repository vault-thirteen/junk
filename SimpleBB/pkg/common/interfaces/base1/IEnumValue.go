package base

import (
	"database/sql/driver"
)

type IEnumValue interface {
	Scan(src any) error
	Value() (driver.Value, error)
	ToString() string
	RawValue() byte
	AsByte() byte
	AsInt() int
	UnmarshalJSON(data []byte) error
	MarshalJSON() ([]byte, error)
}
