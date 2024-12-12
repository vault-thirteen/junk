package s

import (
	"errors"
)

const (
	ErrHttpHost = "HTTP host is not set"
	ErrHttpPort = "HTTP port is not set"
)

// HttpSettings are settings of a generic HTTP server.
type HttpSettings struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func (hs HttpSettings) Check() (err error) {
	if len(hs.Host) == 0 {
		return errors.New(ErrHttpHost)
	}
	if hs.Port == 0 {
		return errors.New(ErrHttpPort)
	}

	return nil
}
