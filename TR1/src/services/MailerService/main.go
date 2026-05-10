package main

import (
	"github.com/vault-thirteen/TR1/src/components/ConsoleComponent"
	"github.com/vault-thirteen/TR1/src/components/ErrorListenerComponent"
	"github.com/vault-thirteen/TR1/src/components/MailerComponent"
	"github.com/vault-thirteen/TR1/src/components/RpcServerComponent"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/services/MailerService/controller"
)

func main() {
	// Order of components must be synchronised with a list of component
	// indices of the controller.
	var serviceComponents = []interfaces.IServiceComponent{
		&coc.ConsoleComponent{},
		&elc.ErrorListenerComponent{},
		&mc.MailerComponent{},
		&rsc.RpcServerComponent{},
	}

	var controller interfaces.IController
	controller = c.NewController()

	app, err := cm.NewApplication(rm.ServiceName_MailerService, serviceComponents, controller)
	mustBeNoError(err)

	err = app.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
