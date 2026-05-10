package c

import (
	"log"

	"github.com/vault-thirteen/JSON-RPC-M1"
	rcs "github.com/vault-thirteen/RingCaptcha/server"
	"github.com/vault-thirteen/TR1/src/components/CaptchaComponent"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/libraries/scheduler"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

// List of component indices of the controller must be synchronised with the
// order of components used in the application's constructor.
const (
	ComponentIndex_ConsoleComponent       = 0
	ComponentIndex_ErrorListenerComponent = 1
	ComponentIndex_CaptchaComponent       = 2
	ComponentIndex_RpcServerComponent     = 3
)

type Controller struct {
	cfg        *cm.Configuration
	errorsChan *chan error
	service    *cm.Service
	far        ControllerFastAccessRegistry
}

func NewController() (c *Controller) {
	errorsChan := make(chan error, 1)

	return &Controller{
		errorsChan: &errorsChan,
	}
}

func (c *Controller) GetRpcFunctions() []jrm1.RpcFunction {
	return []jrm1.RpcFunction{
		c.Ping,
		c.CreateCaptcha,
		c.CheckCaptcha,
		c.HasCaptcha,
	}
}

func (c *Controller) GetScheduledFunctions() []sch.ScheduledFn {
	return []sch.ScheduledFn{}
}

func (c *Controller) GetErrorsChan() (errorsChan *chan error) {
	return c.errorsChan
}

func (c *Controller) LinkWithService(service interfaces.IService) (err error) {
	c.cfg = (service.GetConfiguration()).(*cm.Configuration)
	c.service = service.(*cm.Service)
	c.initFAR()

	return nil
}

func (c *Controller) initFAR() {
	c.far = ControllerFastAccessRegistry{}
	c.far.systemSettings = c.cfg.GetComponent(cm.Component_System, cm.Protocol_None)

	c.far.cc = cc.FromAny(c.service.GetComponentByIndex(ComponentIndex_CaptchaComponent))
	c.far.cs = c.far.cc.GetCaptchaServer()
}

func (c *Controller) GetCaptchaServer() *rcs.CaptchaServer {
	return c.far.cs
}

func (c *Controller) logError(err error) {
	if err == nil {
		return
	}

	if c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode) {
		log.Println(err)
	}
}
