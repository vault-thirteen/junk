package c

import (
	"github.com/vault-thirteen/TR1/src/components/DatabaseComponent"
	"github.com/vault-thirteen/TR1/src/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/models/rpc/Client"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
	"gorm.io/gorm"
)

type ControllerFastAccessRegistry struct {
	systemSettings  *ccse.CommonConfigurationServiceEntry
	messageSettings *ccse.CommonConfigurationServiceEntry

	authServiceClient *rmc.Client

	pageSize        int
	messageEditTime int

	dbc *dc.DatabaseComponent
	db  *gorm.DB
	rcc *rcc.RpcClientComponent
}
