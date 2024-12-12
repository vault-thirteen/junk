package settings

import (
	"errors"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

// SmtpSettings are parameters of the SMTP server used for sending e-mail
// messages. When a password is not set, it is taken from the stdin.
type SmtpSettings struct {
	Host      string `json:"host"`
	Port      uint16 `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"`
}

func (s SmtpSettings) Check() (err error) {
	if (len(s.Host) == 0) ||
		(s.Port == 0) ||
		(len(s.User) == 0) ||
		(len(s.UserAgent) == 0) {
		return errors.New(c.MsgSmtpSettingError)
	}

	return nil
}
