package jmc

import (
	"fmt"

	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/libraries/km"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

type JwtManagerComponent struct {
	cfg   interfaces.IConfiguration
	jwtkm *km.KeyMaker
}

func (c *JwtManagerComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	jwtSettings := cfg.GetComponent(cm.Component_Jwt, cm.Protocol_None)

	kms := &km.KeyMakerSettings{
		SigningMethodName:  jwtSettings.GetParameterAsString(ccp.SigningMethod),
		PrivateKeyFilePath: jwtSettings.GetParameterAsString(ccp.PrivateKeyFilePath),
		PublicKeyFilePath:  jwtSettings.GetParameterAsString(ccp.PublicKeyFilePath),
		IsCacheEnabled:     jwtSettings.GetParameterAsBool(ccp.IsCacheEnabled),
		CacheSizeLimit:     jwtSettings.GetParameterAsInt(ccp.CacheSizeLimit),
		CacheRecordTtl:     uint(jwtSettings.GetParameterAsInt(ccp.CacheRecordTtl)),
	}

	var jwtkm *km.KeyMaker
	jwtkm, err = km.New(kms)
	if err != nil {
		return nil, err
	}

	jmc := &JwtManagerComponent{
		cfg:   cfg,
		jwtkm: jwtkm,
	}

	return jmc, nil
}
func (c *JwtManagerComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *JwtManagerComponent) Start(s interfaces.IService) (err error) {
	return nil
}
func (c *JwtManagerComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	c.ReportStop()

	return nil
}

func (c *JwtManagerComponent) ReportStart() {
	fmt.Println("JwtManagerComponent has started")
}
func (c *JwtManagerComponent) ReportStop() {
	fmt.Println("JwtManagerComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *JwtManagerComponent) {
	return x.(*JwtManagerComponent)
}

// Non-standard methods.

func (c *JwtManagerComponent) GetKeyMaker() (jwtkm *km.KeyMaker) {
	return c.jwtkm
}
