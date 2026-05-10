package rcc

import (
	"fmt"

	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/models/rpc/Client"
	"github.com/vault-thirteen/TR1/src/models/rpc/Proxy"
)

type RpcClientComponent struct {
	cfg       interfaces.IConfiguration
	clientMap map[string]*rmc.Client
	proxyMap  map[string]*rmp.Proxy
}

func (c *RpcClientComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	rcc := &RpcClientComponent{
		cfg:       cfg,
		clientMap: make(map[string]*rmc.Client),
		proxyMap:  make(map[string]*rmp.Proxy),
	}

	err = rcc.addClientIfNeeded(cfg, cm.ClientType_Auth, cm.Protocol_HTTPS, rm.ServiceShortName_Auth)
	if err != nil {
		return nil, err
	}

	err = rcc.addProxyIfNeeded(cfg, cm.ClientType_Captcha, cm.Protocol_HTTP, rm.ServiceShortName_Captcha)
	if err != nil {
		return nil, err
	}

	err = rcc.addClientIfNeeded(cfg, cm.ClientType_Mailer, cm.Protocol_HTTPS, rm.ServiceShortName_Mailer)
	if err != nil {
		return nil, err
	}

	err = rcc.addClientIfNeeded(cfg, cm.ClientType_Message, cm.Protocol_HTTPS, rm.ServiceShortName_Message)
	if err != nil {
		return nil, err
	}

	err = rcc.addClientIfNeeded(cfg, cm.ClientType_RCS, cm.Protocol_HTTPS, rm.ServiceShortName_RCS)
	if err != nil {
		return nil, err
	}

	return rcc, nil
}
func (c *RpcClientComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *RpcClientComponent) addClientIfNeeded(cfg interfaces.IConfiguration, cType string, cProtocol string, shortName string) (err error) {
	clientSettings := cfg.GetClient(cType, cProtocol)

	if clientSettings == nil {
		return nil
	}

	var client *rmc.Client
	client, err = rmc.NewClientFromSettings(clientSettings, shortName)
	if err != nil {
		return err
	}

	c.clientMap[shortName] = client

	return nil
}
func (c *RpcClientComponent) addProxyIfNeeded(cfg interfaces.IConfiguration, cType string, cProtocol string, shortName string) (err error) {
	proxySettings := cfg.GetClient(cType, cProtocol)

	if proxySettings == nil {
		return nil
	}

	var proxy *rmp.Proxy
	proxy, err = rmp.NewProxyFromSettings(proxySettings, shortName)
	if err != nil {
		return err
	}

	c.proxyMap[shortName] = proxy

	return nil
}

func (c *RpcClientComponent) Start(s interfaces.IService) (err error) {
	return nil
}
func (c *RpcClientComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	c.ReportStop()

	return nil
}

func (c *RpcClientComponent) ReportStart() {
	fmt.Println("RpcClientComponent has started")
}
func (c *RpcClientComponent) ReportStop() {
	fmt.Println("RpcClientComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *RpcClientComponent) {
	return x.(*RpcClientComponent)
}

// Non-standard methods.

func (c *RpcClientComponent) GetClientMap() (clientMap map[string]*rmc.Client) {
	return c.clientMap
}
func (c *RpcClientComponent) GetProxyMap() (proxyMap map[string]*rmp.Proxy) {
	return c.proxyMap
}
