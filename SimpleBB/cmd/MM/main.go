package main

import (
	m "github.com/vault-thirteen/SimpleBB/pkg/MM/server"
	app2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"log"

	ms "github.com/vault-thirteen/SimpleBB/pkg/MM/settings"
)

func main() {
	theApp, err := app2.NewApplication[*ms.Settings, *m.Server](&ms.Settings{}, &m.Server{}, app2.ServiceName_MM, app2.ConfigurationFilePathDefault_MM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
