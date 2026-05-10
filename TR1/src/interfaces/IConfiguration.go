package interfaces

import (
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
)

type IConfiguration interface {
	GetComponent(cType string, cProtocol string) *ccse.CommonConfigurationServiceEntry
	GetServer(cType string, cProtocol string) *ccse.CommonConfigurationServiceEntry
	GetClient(cType string, cProtocol string) *ccse.CommonConfigurationServiceEntry
}
