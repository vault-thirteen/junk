package models

import "flag"

// CliArguments is command line interface arguments.
type CliArguments struct {
	HeaderFilePath                string
	FuncProcessorSettingsFilePath string
}

// NewCliArgumentsFromOsArgs creates CLI arguments from OS arguments.
func NewCliArgumentsFromOsArgs() (ca *CliArguments, err error) {
	flagHeaderFilePath := flag.String("header", "", "path to C header file")
	flagFuncProcessorSettingsFilePath := flag.String("settings", "", "path to settings file")
	flag.Parse()

	ca = &CliArguments{
		HeaderFilePath:                *flagHeaderFilePath,
		FuncProcessorSettingsFilePath: *flagFuncProcessorSettingsFilePath,
	}

	return ca, nil
}
