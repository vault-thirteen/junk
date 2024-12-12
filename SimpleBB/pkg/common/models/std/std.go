package std

import "fmt"

const (
	ByteMin = 0
	ByteMax = 255
)

const (
	ErrF_ByteOverflow = "byte overflow: %v"
	ErrF_BoolOverflow = "bool overflow: %v"
)

func CastIntToByte(i int) (b byte, err error) {
	if (i < ByteMin) || (i > ByteMax) {
		return b, fmt.Errorf(ErrF_ByteOverflow, i)
	}

	return byte(i), nil
}

func CastIntToBool(i int) (b bool, err error) {
	switch i {
	case 0:
		return false, nil

	case 1:
		return true, nil

	default:
		return b, fmt.Errorf(ErrF_BoolOverflow, i)
	}
}
