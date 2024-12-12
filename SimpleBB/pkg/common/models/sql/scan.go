package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/std"
	"reflect"
	"strconv"

	num "github.com/vault-thirteen/auxie/number"
)

const (
	ErrF_CanNotScanSourceAsInt     = "source can not be scanned as int: %s"
	ErrF_CanNotScanSourceAsString  = "source can not be scanned as string: %s"
	ErrF_CanNotScanSourceAsBoolean = "source can not be scanned as bool: %s"
)

func ScanSrcAsInt(src any) (i int, err error) {
	switch src.(type) {
	case int64:
		i = int(src.(int64))

	case []byte:
		i, err = num.ParseInt(string(src.([]byte)))
		if err != nil {
			return 0, err
		}

	case string:
		i, err = num.ParseInt(src.(string))
		if err != nil {
			return 0, err
		}

	default:
		return 0, fmt.Errorf(ErrF_CanNotScanSourceAsInt, reflect.TypeOf(src).String())
	}

	return i, nil
}

func ScanSrcAsByte(src any) (b byte, err error) {
	var i int
	i, err = ScanSrcAsInt(src)
	if err != nil {
		return b, err
	}

	return std.CastIntToByte(i)
}

func ScanSrcAsString(src any) (s string, err error) {
	switch src.(type) {
	case []byte:
		s = string(src.([]byte))

	case string:
		s = src.(string)

	default:
		return s, fmt.Errorf(ErrF_CanNotScanSourceAsString, reflect.TypeOf(src).String())
	}

	return s, nil
}

func ScanSrcAsBoolean(src any) (b bool, err error) {
	// 1. Try the easiest way first.
	var i int
	i, err = ScanSrcAsInt(src)
	if err == nil {
		return std.CastIntToBool(i)
	}

	// 2. Source is not an integer number.
	switch src.(type) {
	case bool:
		b = src.(bool)

	case []byte:
		b, err = strconv.ParseBool(string(src.([]byte)))
		if err != nil {
			return b, err
		}

	case string:
		b, err = strconv.ParseBool(src.(string))
		if err != nil {
			return b, err
		}

	default:
		return b, fmt.Errorf(ErrF_CanNotScanSourceAsBoolean, reflect.TypeOf(src).String())
	}

	return b, nil
}

func NewArrayFromScannableSource[T any](src base.IScannableSequence) (values []T, err error) {
	values = []T{}
	var value T

	for src.Next() {
		err = src.Scan(&value)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

func NewValueFromScannableSource[T any](src base.IScannable) (*T, error) {
	var value = new(T)

	err := src.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return value, nil
}

func NewNonNullValueFromScannableSource[T any](src base.IScannable) (T, error) {
	var value T

	err := src.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return value, nil
		} else {
			return value, err
		}
	}

	return value, nil
}
