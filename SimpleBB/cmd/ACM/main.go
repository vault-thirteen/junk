package main

import (
	a "github.com/vault-thirteen/SimpleBB/pkg/ACM/server"
	app2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"log"

	as "github.com/vault-thirteen/SimpleBB/pkg/ACM/settings"
)

func main() {
	theApp, err := app2.NewApplication[*as.Settings, *a.Server](&as.Settings{}, &a.Server{}, app2.ServiceName_ACM, app2.ConfigurationFilePathDefault_ACM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
