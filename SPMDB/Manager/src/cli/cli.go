package cli

import (
	"errors"
	"flag"
	"fmt"
	"spdb_manager/mode"
)

type Parameters struct {
	Mode             string
	DbFile           string
	VideosListUrl    string
	ScriptsFolder    string
	ServerHost       string
	ServerPort       uint16
	FilesToSkip      uint
	GuiScriptsFolder string
}

// GetParameters get the command line arguments which were used for starting
// this application.
func GetParameters() (params *Parameters, err error) {
	cliMode := flag.String("mode", "", "work mode")
	cliDbFile := flag.String("database", "", "path to a database file")
	cliVideosListUrl := flag.String("videos", "", "URL of a file with a list of videos")
	cliScriptsFolder := flag.String("scripts", "", "scripts folder")
	cliServerHost := flag.String("host", "localhost", "server host")
	cliServerPort := flag.Uint("port", 2000, "server port")
	cliSkipFiles := flag.Uint("skip", 0, "number of files to skip at the beginning of import")
	cliGuiScriptsFolder := flag.String("gui", `gui`, "gui scripts folder")
	flag.Parse()

	params = &Parameters{
		Mode:             *cliMode,
		DbFile:           *cliDbFile,
		VideosListUrl:    *cliVideosListUrl,
		ScriptsFolder:    *cliScriptsFolder,
		ServerHost:       *cliServerHost,
		ServerPort:       uint16(*cliServerPort),
		FilesToSkip:      *cliSkipFiles,
		GuiScriptsFolder: *cliGuiScriptsFolder,
	}

	// Here we check only the common parameters.
	// Specific parameters are checked in each specific handler.

	if !mode.IsAvailable(params.Mode) {
		return nil, fmt.Errorf("unsupported mode: '%s'", params.Mode)
	}

	if len(params.DbFile) == 0 {
		return nil, errors.New("database file is not set")
	}

	return params, nil
}
