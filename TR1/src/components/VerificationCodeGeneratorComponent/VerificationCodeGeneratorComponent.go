package vcgc

import (
	"fmt"

	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	rp "github.com/vault-thirteen/auxie/rpofs"
)

type VerificationCodeGeneratorComponent struct {
	cfg interfaces.IConfiguration
	vcg *rp.Generator
}

func (c *VerificationCodeGeneratorComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	symbols := cm.MakeSymbolsNumbersAndCapitalLatinLetters()

	systemSettings := cfg.GetComponent(cm.Component_System, cm.Protocol_None)

	var vcg *rp.Generator
	vcg, err = rp.NewGenerator(systemSettings.GetParameterAsInt(ccp.VerificationCodeLength), symbols)
	if err != nil {
		return nil, err
	}

	jmc := &VerificationCodeGeneratorComponent{
		cfg: cfg,
		vcg: vcg,
	}

	return jmc, nil
}
func (c *VerificationCodeGeneratorComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *VerificationCodeGeneratorComponent) Start(s interfaces.IService) (err error) {
	return nil
}
func (c *VerificationCodeGeneratorComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	c.ReportStop()

	return nil
}

func (c *VerificationCodeGeneratorComponent) ReportStart() {
	fmt.Println("VerificationCodeGeneratorComponent has started")
}
func (c *VerificationCodeGeneratorComponent) ReportStop() {
	fmt.Println("VerificationCodeGeneratorComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *VerificationCodeGeneratorComponent) {
	return x.(*VerificationCodeGeneratorComponent)
}

// Non-standard methods.

func (c *VerificationCodeGeneratorComponent) GetVcg() (vcg *rp.Generator) { return c.vcg }
