package json

import (
	"encoding/json"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/std"
)

func UnmarshalAsInt(src []byte) (i int, err error) {
	err = json.Unmarshal(src, &i)
	if err != nil {
		return i, err
	}

	return i, nil
}

func UnmarshalAsByte(src []byte) (b byte, err error) {
	var i int
	i, err = UnmarshalAsInt(src)
	if err != nil {
		return b, err
	}

	return std.CastIntToByte(i)
}

func UnmarshalAsString(src []byte) (s string, err error) {
	err = json.Unmarshal(src, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func UnmarshalAsBoolean(src []byte) (b bool, err error) {
	err = json.Unmarshal(src, &b)
	if err != nil {
		return b, err
	}

	return b, nil
}
