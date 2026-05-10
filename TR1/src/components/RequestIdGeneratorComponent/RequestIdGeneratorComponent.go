package rigc

import (
	"fmt"

	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	rp "github.com/vault-thirteen/auxie/rpofs"
)

type RequestIdGeneratorComponent struct {
	cfg  interfaces.IConfiguration
	ridg *rp.Generator
}

func (c *RequestIdGeneratorComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	symbols := cm.MakeSymbolsNumbersAndCapitalLatinLetters()

	systemSettings := cfg.GetComponent(cm.Component_System, cm.Protocol_None)

	var ridg *rp.Generator
	ridg, err = rp.NewGenerator(systemSettings.GetParameterAsInt(ccp.RequestIdLength), symbols)
	if err != nil {
		return nil, err
	}

	jmc := &RequestIdGeneratorComponent{
		cfg:  cfg,
		ridg: ridg,
	}

	return jmc, nil
}
func (c *RequestIdGeneratorComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *RequestIdGeneratorComponent) Start(s interfaces.IService) (err error) {
	return nil
}
func (c *RequestIdGeneratorComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	c.ReportStop()

	return nil
}

func (c *RequestIdGeneratorComponent) ReportStart() {
	fmt.Println("RequestIdGeneratorComponent has started")
}
func (c *RequestIdGeneratorComponent) ReportStop() {
	fmt.Println("RequestIdGeneratorComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *RequestIdGeneratorComponent) {
	return x.(*RequestIdGeneratorComponent)
}

// Non-standard methods.

func (c *RequestIdGeneratorComponent) GetRidg() (ridg *rp.Generator) {
	return c.ridg
}
