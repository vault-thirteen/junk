package base

import (
	"database/sql/driver"
)

type IEnum interface {
	SetValue(value IEnumValue) error
	SetValueFast(value IEnumValue)
	GetValue() IEnumValue
	Scan(src any) error
	Value() (driver.Value, error)
	ToString() string
	AsByte() byte
	AsInt() int
	UnmarshalJSON(data []byte) error
	MarshalJSON() ([]byte, error)
}
