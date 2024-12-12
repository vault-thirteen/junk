package main

import (
	g "github.com/vault-thirteen/SimpleBB/pkg/GWM/server"
	app2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"log"

	gs "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
)

func main() {
	theApp, err := app2.NewApplication[gs.ISettings, *g.Server](gs.NewSettings(), &g.Server{}, app2.ServiceName_GWM, app2.ConfigurationFilePathDefault_GWM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
