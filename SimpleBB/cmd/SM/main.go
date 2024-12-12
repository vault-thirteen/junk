package main

import (
	s "github.com/vault-thirteen/SimpleBB/pkg/SM/server"
	app2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"log"

	ss "github.com/vault-thirteen/SimpleBB/pkg/SM/settings"
)

func main() {
	theApp, err := app2.NewApplication[*ss.Settings, *s.Server](&ss.Settings{}, &s.Server{}, app2.ServiceName_SM, app2.ConfigurationFilePathDefault_SM)
	mustBeNoError(err)

	err = theApp.Use()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
