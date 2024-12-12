package s

import (
	"errors"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

// SystemSettings are system settings.
type SystemSettings struct {
	PageSize    base2.Count `json:"pageSize"`
	DKeySize    base2.Count `json:"dKeySize"`
	IsDebugMode base2.Flag  `json:"isDebugMode"`
}

func (s SystemSettings) Check() (err error) {
	if (s.PageSize == 0) ||
		(s.DKeySize == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	return nil
}
