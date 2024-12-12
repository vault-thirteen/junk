package main

import (
	r "github.com/vault-thirteen/SimpleBB/pkg/RCS/server"
	app2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"log"

	rs "github.com/vault-thirteen/SimpleBB/pkg/RCS/settings"
)

func main() {
	theApp, err := app2.NewApplication[*rs.Settings, *r.Server](&rs.Settings{}, &r.Server{}, app2.ServiceName_RCS, app2.ConfigurationFilePathDefault_RCS)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
