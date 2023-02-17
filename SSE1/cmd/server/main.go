package main

import (
	"flag"
	"log"
	"os"

	"github.com/vault-thirteen/junk/SSE1/pkg/models/app"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/configuration"
)

// Messages.
const (
	MsgCommandLineArguments     = "Command Line Arguments"
	MsgApplicationConfiguration = "Application Configuration"
	MsgApplication              = "Application"
)

// OS Exit Codes.
const (
	OsExitCodeError = 1
)

func main() {
	var err error
	var cla CommandLineArguments
	cla, err = getCommandLineArguments()
	mustBeNoError(MsgCommandLineArguments, err)
	var appConfiguration *configuration.AppConfiguration
	appConfiguration, err = configuration.NewAppConfiguration(cla.PathToConfigurationFile)
	mustBeNoError(MsgApplicationConfiguration, err)
	var application *app.Application
	application, err = app.NewApplication(appConfiguration)
	mustBeNoError(MsgApplication, err)
	err = application.Run()
	checkError(err)
}

func getCommandLineArguments() (cla CommandLineArguments, err error) {
	flag.StringVar(
		&cla.PathToConfigurationFile,
		PathToConfigurationFileArgName,
		PathToConfigurationFileDefaultValue,
		PathToConfigurationFileUsageHint,
	)
	flag.Parse()
	return
}

func mustBeNoError(
	errorPrefix string,
	err error,
) {
	const SymbolColon = ":"
	if err != nil {
		log.Println(errorPrefix+SymbolColon, err)
		os.Exit(OsExitCodeError)
	}
}

func checkError(
	err error,
) {
	if err != nil {
		log.Println(err)
	}
}
