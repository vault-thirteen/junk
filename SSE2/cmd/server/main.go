package main

import (
	"log"

	"github.com/vault-thirteen/junk/SSE2/internal/application"
)

func main() {
	app, err := application.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Start()
	app.MustBeNoError(err)

	err = app.WaitForQuitSignal()
	app.MustBeNoError(err)
}
