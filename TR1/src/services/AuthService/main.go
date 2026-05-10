package main

import (
	"github.com/vault-thirteen/TR1/src/components/ConsoleComponent"
	"github.com/vault-thirteen/TR1/src/components/DatabaseComponent"
	"github.com/vault-thirteen/TR1/src/components/ErrorListenerComponent"
	"github.com/vault-thirteen/TR1/src/components/JwtManagerComponent"
	"github.com/vault-thirteen/TR1/src/components/RequestIdGeneratorComponent"
	"github.com/vault-thirteen/TR1/src/components/RpcClientComponent"
	"github.com/vault-thirteen/TR1/src/components/RpcServerComponent"
	"github.com/vault-thirteen/TR1/src/components/SchedulerComponent"
	"github.com/vault-thirteen/TR1/src/components/VerificationCodeGeneratorComponent"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/models/rpc"
	"github.com/vault-thirteen/TR1/src/services/AuthService/controller"
)

func main() {
	// Order of components must be synchronised with a list of component
	// indices of the controller.
	var serviceComponents = []interfaces.IServiceComponent{
		&coc.ConsoleComponent{},
		&elc.ErrorListenerComponent{},
		&dc.DatabaseComponent{},
		&jmc.JwtManagerComponent{},
		&rigc.RequestIdGeneratorComponent{},
		&rcc.RpcClientComponent{},
		&vcgc.VerificationCodeGeneratorComponent{},
		&rsc.RpcServerComponent{},
		&shc.SchedulerComponent{},
	}

	var controller interfaces.IController
	controller = c.NewController()

	app, err := cm.NewApplication(rm.ServiceName_AuthService, serviceComponents, controller)
	mustBeNoError(err)

	err = app.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
