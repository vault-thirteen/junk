package cm

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
	"golang.org/x/term"
)

const (
	ErrF_FilterTypeIsUnsupported = "filter type is unsupported: %s"
)

const (
	Field_Type     = "type"
	Field_Protocol = "protocol"
)

const (
	Component_Captcha  = "captcha"  // Captcha settings.
	Component_Database = "database" // Database settings.
	Component_SFS      = "sfs"      // Static Files Server settings.
	Component_Jwt      = "jwt"      // JWT settings.
	Component_Message  = "message"  // Message settings.
	Component_Role     = "role"     // Role settings.
	Component_Mailer   = "mailer"   // Mailer settings.
	Component_System   = "system"   // System settings.
)

const (
	Protocol_None  = ""
	Protocol_MySQL = "mysql"
	Protocol_HTTP  = "http"
	Protocol_HTTPS = "https"
)

const (
	ServerType_Internal = "internal"
	ServerType_External = "external"
)

const (
	ClientType_Auth    = "auth"
	ClientType_Captcha = "captcha" // Captcha images (Proxy).
	ClientType_Mailer  = "mailer"
	ClientType_Message = "message"
	ClientType_RCS     = "rcs" // Captcha questions (RPC).
)

type Configuration struct {
	Service ConfigurationService
}

func NewConfigurationFromFile(filePath string) (c *Configuration, err error) {
	var buf []byte
	buf, err = os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	cr := &ConfigurationRaw{}
	err = json.Unmarshal(buf, cr)
	if err != nil {
		return nil, err
	}

	c = &Configuration{
		Service: ConfigurationService{
			ShortName: cr.Service.ShortName,
			FullName:  cr.Service.FullName,
		},
	}

	c.Service.Components, err = parseConfigurationServiceEntries(cr.Service.Components)
	if err != nil {
		return nil, err
	}

	c.Service.Servers, err = parseConfigurationServiceEntries(cr.Service.Servers)
	if err != nil {
		return nil, err
	}

	c.Service.Clients, err = parseConfigurationServiceEntries(cr.Service.Clients)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c Configuration) GetComponent(cType string, cProtocol string) *ccse.CommonConfigurationServiceEntry {
	return getConfigurationServiceEntryByTypeAndProtocol(c.Service.Components, cType, cProtocol)
}
func (c Configuration) GetServer(cType string, cProtocol string) *ccse.CommonConfigurationServiceEntry {
	return getConfigurationServiceEntryByTypeAndProtocol(c.Service.Servers, cType, cProtocol)
}
func (c Configuration) GetClient(cType string, cProtocol string) *ccse.CommonConfigurationServiceEntry {
	return getConfigurationServiceEntryByTypeAndProtocol(c.Service.Clients, cType, cProtocol)
}

func GetPasswordFromStdin(object string) (pwd string, err error) {
	msg := fmt.Sprintf("Enter password for %s:", object)
	fmt.Println(msg)

	var buf []byte
	buf, err = term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
