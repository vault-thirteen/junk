package main

import (
	"os"
	"spdb_manager/cli"
	"spdb_manager/database"
)

// createEmptyDb creates an empty database.
func createEmptyDb(
	cliParams *cli.Parameters,
) (err error) {
	err = createEmptyFile(cliParams.DbFile)
	if err != nil {
		return err
	}

	err = database.CheckConnection(cliParams.DbFile)
	if err != nil {
		return err
	}

	return nil
}

// createEmptyFile creates an empty file.
func createEmptyFile(filePath string) (err error) {
	var file *os.File
	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
