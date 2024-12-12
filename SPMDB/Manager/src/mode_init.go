package main

import (
	"errors"
	"path/filepath"
	"spdb_manager/cli"
	"spdb_manager/cli/printer"
	"spdb_manager/database"
	"spdb_manager/helper"

	"go.uber.org/multierr"
)

// initDb runs SQL scripts in an empty database.
func initDb(
	cliParams *cli.Parameters,
) (err error) {
	if len(cliParams.ScriptsFolder) == 0 {
		return errors.New("scripts folder is not set")
	}

	var files []string
	files, err = helper.GetDirectoryFilesList(cliParams.ScriptsFolder)
	if err != nil {
		return err
	}

	var dbo *database.DbSimple
	dbo, err = database.OpenSimple(cliParams.DbFile)
	if err != nil {
		return err
	}

	defer func() {
		derr := dbo.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	for _, file := range files {
		filePath := filepath.Join(cliParams.ScriptsFolder, file)
		printer.PrintInfo(filePath)

		err = dbo.RunScript(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}
