package main

import (
	"fmt"
	"os"
	"spdb_manager/cli"
	"spdb_manager/cli/printer"
	"spdb_manager/mode"
	"spdb_manager/win32"
)

// main is an application's entry point.
func main() {
	var err error

	err = win32.EnableVirtualTerminalProcessing()
	mustBeNoError(err)

	var cliParams *cli.Parameters
	cliParams, err = cli.GetParameters()
	if err != nil {
		printer.PrintError(fmt.Sprintf("Error: %s.", err.Error()))
		printer.PrintNormalText("Use the help ('-h' parameter) to view the list of available parameters.")
		os.Exit(1)
	}

	quitChan := make(chan bool)

	switch cliParams.Mode {
	case mode.WorkModeCreateDb:
		err = createEmptyDb(cliParams)
	case mode.WorkModeInitDb:
		err = initDb(cliParams)
	case mode.WorkModeImportVideosToDb:
		err = importNewMoviesList(cliParams)
	case mode.WorkModeEditDb:
		err = editData(cliParams, &quitChan)
	case mode.WorkModeViewDb:
		err = viewData(cliParams, &quitChan)
	}
	mustBeNoError(err)

	switch cliParams.Mode {
	case mode.WorkModeEditDb,
		mode.WorkModeViewDb:
		<-quitChan
	}
}

// mustBeNoError closes the application on error.
func mustBeNoError(err error) {
	if err != nil {
		printer.PrintError(err.Error())
		os.Exit(1)
	}
}

//TODO:Fix names in DB.
