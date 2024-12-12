package s

import (
	"errors"
	"fmt"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

const (
	ErrScsSchema = "SCS schema is not set"
	ErrScsHost   = "SCS host is not set"
	ErrScsPort   = "SCS port is not set"
	ErrScsPath   = "SCS path is not set"
)

// ServiceClientSettings are common settings for a service client.
type ServiceClientSettings struct {
	Schema                      string `json:"schema"`
	Host                        string `json:"host"`
	Port                        uint16 `json:"port"`
	Path                        string `json:"path"`
	EnableSelfSignedCertificate bool   `json:"enableSelfSignedCertificate"`
}

// Check checks common settings of a service client.
func (scs ServiceClientSettings) Check() (err error) {
	if len(scs.Schema) == 0 {
		return errors.New(ErrScsSchema)
	}
	if len(scs.Host) == 0 {
		return errors.New(ErrScsHost)
	}
	if scs.Port == 0 {
		return errors.New(ErrScsPort)
	}
	if len(scs.Path) == 0 {
		return errors.New(ErrScsPath)
	}

	return nil
}

func DetailedScsError(serviceShortName string, errIn error) (errOut error) {
	return fmt.Errorf(c.MsgFServiceClientSettingsError, serviceShortName, errIn.Error())
}
