package settings

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// SystemSettings are system settings.
type SystemSettings struct {
	IsDebugMode cmb.Flag `json:"isDebugMode"`
}

func (s SystemSettings) Check() (err error) {
	return nil
}
