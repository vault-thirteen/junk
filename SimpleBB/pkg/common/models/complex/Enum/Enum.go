package enum

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	ev "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
)

// Enum is a simple enumeration type supporting up to 255 distinct values.
// It is not recommended to use zero as a value. Zero is commonly used as a
// marker of a non-set value. The minimum value should always be equal to 1.

// Limitations for enumeration value.
const (
	// EnumValueMin is a minimal value of an enumeration.
	EnumValueMin = 1

	// Maximum value is set in a child class. Oops. Golang does not know what
	// classes are. So, we emulate some parts of the class behaviour as we can.
)

const (
	Err_MaxValue           = "max value error"
	ErrF_EnumValueOverflow = "enumeration value overflow: %v"
)

//TODO: Re-visit this page in 10 years. It may be so that, in 10 years from now
// Go language will introduce classes or objects with constructors.

// This model must be initialised manually using a constructor while Go
// language does not have constructors. Go language does not have constructors
// as there are no classes in the language. So-called "structs" in Go language
// do not have an obligatory constructor. All this makes the process of object
// creation practically uncontrollable. To fix all the automatic
// initialisations in the code, all the code must be re-written. The chain of
// dependencies will overwhelm you if you decide to re-write all the parts
// where objects are created. This is why all serious programming languages
// use built-in constructors for objects which can not be fooled as in Golang.
// All in all, Go language is a complete shit.

type enum struct {
	value    base.IEnumValue
	maxValue base.IEnumValue
}

func NewEnum(maxValue base.IEnumValue) (base.IEnum, error) {
	e := newEnum(maxValue)

	err := e.checkMaxValue(maxValue)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func NewEnumFast(maxValue base.IEnumValue) base.IEnum {
	e := newEnum(maxValue)

	err := e.checkMaxValue(maxValue)
	if err != nil {
		panic(err)
	}

	return e
}

func newEnum(maxValue base.IEnumValue) *enum {
	return &enum{
		value:    ev.NewEnumValue(0),
		maxValue: maxValue,
	}
}

func (e *enum) checkMaxValue(maxValue base.IEnumValue) (err error) {
	if maxValue.RawValue() < EnumValueMin {
		return errors.New(Err_MaxValue)
	}
	return nil
}

func (e *enum) SetValue(value base.IEnumValue) (err error) {
	err = e.checkValue(value)
	if err != nil {
		return err
	}

	e.value = value
	return nil
}

func (e *enum) SetValueFast(value base.IEnumValue) {
	err := e.checkValue(value)
	if err != nil {
		panic(err)
	}

	e.value = value
}

func (e *enum) checkValue(v base.IEnumValue) (err error) {
	if (v.RawValue() < EnumValueMin) || (v.RawValue() > e.maxValue.RawValue()) {
		return fmt.Errorf(ErrF_EnumValueOverflow, v)
	}
	return nil
}

func (e *enum) GetValue() base.IEnumValue {
	return e.value
}

func (e *enum) Scan(src any) (err error) {
	err = e.value.Scan(src)
	if err != nil {
		return err
	}

	return e.checkValue(e.value)
}

func (e enum) Value() (dv driver.Value, err error) {
	return e.value.Value()
}

func (e enum) ToString() string {
	return e.value.ToString()
}

func (e enum) AsByte() byte {
	return e.value.AsByte()
}

func (e enum) AsInt() int {
	return e.value.AsInt()
}

func (e *enum) UnmarshalJSON(data []byte) (err error) {
	err = e.value.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	return e.checkValue(e.value)
}

func (e enum) MarshalJSON() (ba []byte, err error) {
	return e.value.MarshalJSON()
}
