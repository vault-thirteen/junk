package s

import (
	"errors"
	"fmt"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"os"

	"golang.org/x/term"
)

const (
	ErrSettingsFileIsNotSet = "settings file is not set"
)

const (
	RpcDurationFieldName  = "dur"
	RpcRequestIdFieldName = "rid"
)

func GetPasswordFromStdin(hint string) (pwd string, err error) {
	fmt.Println(hint)

	var buf []byte
	buf, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func CheckSettingsFilePath(sfp cm.Path) (err error) {
	if len(sfp) == 0 {
		return errors.New(ErrSettingsFileIsNotSet)
	}

	return nil
}
