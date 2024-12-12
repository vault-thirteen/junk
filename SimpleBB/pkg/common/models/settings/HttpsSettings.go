package s

import (
	"errors"
)

const (
	ErrHttpsHost     = "HTTPS host is not set"
	ErrHttpsPort     = "HTTPS port is not set"
	ErrHttpsCertFile = "HTTPS CertFile is not set"
	ErrHttpsKeyFile  = "HTTPS KeyFile is not set"
)

// HttpsSettings are settings of a generic HTTPS server.
type HttpsSettings struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

func (hss HttpsSettings) Check() (err error) {
	if len(hss.Host) == 0 {
		return errors.New(ErrHttpsHost)
	}
	if hss.Port == 0 {
		return errors.New(ErrHttpsPort)
	}
	if len(hss.CertFile) == 0 {
		return errors.New(ErrHttpsCertFile)
	}
	if len(hss.KeyFile) == 0 {
		return errors.New(ErrHttpsKeyFile)
	}

	return nil
}
