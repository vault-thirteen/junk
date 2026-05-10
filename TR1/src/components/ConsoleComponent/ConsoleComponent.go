package coc

import (
	"fmt"

	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/common/win32"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

type ConsoleComponent struct {
	cfg interfaces.IConfiguration
}

func (c *ConsoleComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	cc := &ConsoleComponent{
		cfg: cfg,
	}

	systemSettings := cfg.GetComponent(cm.Component_System, cm.Protocol_None)
	shouldInitColours := systemSettings.GetParameterAsBool(ccp.InitConsoleColours)
	if shouldInitColours {
		err = win32.Enable_console_colours()
		if err != nil {
			return nil, err
		}
	}

	return cc, nil
}
func (c *ConsoleComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *ConsoleComponent) Start(s interfaces.IService) (err error) {
	return nil
}
func (c *ConsoleComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	c.ReportStop()

	return nil
}

func (c *ConsoleComponent) ReportStart() {
	fmt.Println("ConsoleComponent has started")
}
func (c *ConsoleComponent) ReportStop() {
	fmt.Println("ConsoleComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *ConsoleComponent) {
	return x.(*ConsoleComponent)
}

// Non-standard methods.
