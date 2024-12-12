package s

import (
	"errors"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

// CaptchaImageSettings are settings of captcha images.
type CaptchaImageSettings struct {
	// Images Server.
	Schema string `json:"schema"`
	Host   string `json:"host"`
	Port   uint16 `json:"port"`
	Path   string `json:"path"`
}

func (s CaptchaImageSettings) Check() (err error) {
	if (len(s.Schema) == 0) ||
		(len(s.Host) == 0) ||
		(s.Port == 0) ||
		(len(s.Path) == 0) {
		return errors.New(c.MsgCaptchaImageServerSettingError)
	}

	return nil
}
