package main

import (
	"spdb_manager/cli"
	"spdb_manager/server"
)

// viewData starts the server in viewer mode.
func viewData(
	cliParams *cli.Parameters,
	quitChan *chan bool,
) (err error) {
	var srv *server.Server
	srv, err = server.New(
		quitChan,
		cliParams.ServerHost,
		cliParams.ServerPort,
		false,
		cliParams.DbFile,
		cliParams.GuiScriptsFolder,
	)
	if err != nil {
		return err
	}

	srv.Start()

	return nil
}
