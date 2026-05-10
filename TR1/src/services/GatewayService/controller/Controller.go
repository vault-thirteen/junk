package c

import (
	"log"

	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/components/HttpServerComponent"
	"github.com/vault-thirteen/TR1/src/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/components/StaticFileServerComponent"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/libraries/scheduler"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

// List of component indices of the controller must be synchronised with the
// order of components used in the application's constructor.
const (
	ComponentIndex_ConsoleComponent          = 0
	ComponentIndex_ErrorListenerComponent    = 1
	ComponentIndex_RpcClientComponent        = 2
	ComponentIndex_HttpServerComponent       = 3
	ComponentIndex_StaticFileServerComponent = 4
)

type Controller struct {
	cfg        *cm.Configuration
	errorsChan *chan error
	service    *cm.Service
	far        ControllerFastAccessRegistry

	apiHandlers                   map[string]rm.RequestHandler
	httpStatusCodesByRpcErrorCode rm.HttpStatusCodeByRpcErrorCode
	publicSettingsFile            []byte
}

func NewController() (c *Controller) {
	errorsChan := make(chan error, 1)

	return &Controller{
		errorsChan: &errorsChan,
	}
}

func (c *Controller) GetRpcFunctions() []jrm1.RpcFunction {
	return []jrm1.RpcFunction{}
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
	c.initAPI()

	err = c.initPublicSettings()
	if err != nil {
		return err
	}

	c.initGatewayRouter()

	return nil
}

func (c *Controller) initFAR() {
	c.far = ControllerFastAccessRegistry{}

	c.far.systemSettings = c.cfg.GetComponent(cm.Component_System, cm.Protocol_None)

	c.far.rcc = rcc.FromAny(c.service.GetComponentByIndex(ComponentIndex_RpcClientComponent))
	c.far.httpServer = hsc.FromAny(c.service.GetComponentByIndex(ComponentIndex_HttpServerComponent))
	c.far.fileServer = sfsc.FromAny(c.service.GetComponentByIndex(ComponentIndex_StaticFileServerComponent)).GetFileServer()

	c.far.authServiceClient = c.far.rcc.GetClientMap()[rm.ServiceShortName_Auth]
	c.far.messageServiceClient = c.far.rcc.GetClientMap()[rm.ServiceShortName_Message]
	c.far.captchaServiceProxy = c.far.rcc.GetProxyMap()[rm.ServiceShortName_Captcha]

	c.far.pageSize = c.far.systemSettings.GetParameterAsInt(ccp.PageSize)
	c.far.messageEditTime = c.far.systemSettings.GetParameterAsInt(ccp.MessageEditTime)
	c.far.isDeveloperMode = c.far.systemSettings.GetParameterAsBool(ccp.IsDeveloperMode)
	c.far.devModeHttpHeaderAccessControlAllowOrigin = c.far.systemSettings.GetParameterAsString(ccp.DeveloperMode_HttpHeader_AccessControlAllowOrigin)
	c.far.clientIPAddressSource_CustomHeader = c.far.systemSettings.GetParameterAsString(ccp.ClientIPAddressSource_CustomHeader)
	c.far.sessionMaxDuration = c.far.systemSettings.GetParameterAsInt(ccp.SessionMaxDuration)
	c.far.cacheControlMaxAge = c.far.systemSettings.GetParameterAsInt(ccp.CacheControlMaxAge)
}

func (c *Controller) logError(err error) {
	if err == nil {
		return
	}

	if c.far.systemSettings.GetParameterAsBool(ccp.IsDebugMode) {
		log.Println(err)
	}
}
