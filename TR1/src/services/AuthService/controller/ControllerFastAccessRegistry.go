package c

import (
	"github.com/vault-thirteen/TR1/src/components/DatabaseComponent"
	"github.com/vault-thirteen/TR1/src/components/JwtManagerComponent"
	"github.com/vault-thirteen/TR1/src/components/RequestIdGeneratorComponent"
	"github.com/vault-thirteen/TR1/src/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/components/VerificationCodeGeneratorComponent"
	"github.com/vault-thirteen/TR1/src/libraries/km"
	"github.com/vault-thirteen/TR1/src/models/rpc/Client"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
	rp "github.com/vault-thirteen/auxie/rpofs"
	"gorm.io/gorm"
)

type ControllerFastAccessRegistry struct {
	systemSettings  *ccse.CommonConfigurationServiceEntry
	messageSettings *ccse.CommonConfigurationServiceEntry
	roleSettings    *ccse.CommonConfigurationServiceEntry

	rcsServiceClient    *rmc.Client
	mailerServiceClient *rmc.Client

	pageSize int

	dbc   *dc.DatabaseComponent
	db    *gorm.DB
	ridgc *rigc.RequestIdGeneratorComponent
	ridg  *rp.Generator
	rcc   *rcc.RpcClientComponent
	vcgc  *vcgc.VerificationCodeGeneratorComponent
	vcg   *rp.Generator
	jmc   *jmc.JwtManagerComponent
	jwtkm *km.KeyMaker
}
