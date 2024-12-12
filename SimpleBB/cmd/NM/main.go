package main

import (
	n "github.com/vault-thirteen/SimpleBB/pkg/NM/server"
	app2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"log"

	ns "github.com/vault-thirteen/SimpleBB/pkg/NM/settings"
)

func main() {
	theApp, err := app2.NewApplication[*ns.Settings, *n.Server](&ns.Settings{}, &n.Server{}, app2.ServiceName_NM, app2.ConfigurationFilePathDefault_NM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
