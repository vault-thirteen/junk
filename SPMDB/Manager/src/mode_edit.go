package main

import (
	"spdb_manager/cli"
	"spdb_manager/server"
)

// editData starts the server in editor mode.
func editData(
	cliParams *cli.Parameters,
	quitChan *chan bool,
) (err error) {
	var srv *server.Server
	srv, err = server.New(
		quitChan,
		cliParams.ServerHost,
		cliParams.ServerPort,
		true,
		cliParams.DbFile,
		cliParams.GuiScriptsFolder,
	)
	if err != nil {
		return err
	}

	srv.Start()

	return nil
}
