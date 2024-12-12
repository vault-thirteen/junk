package main

import (
	"fmt"
	"io"
	"os"

	"github.com/vault-thirteen/SDLW/development/pp"
	"github.com/vault-thirteen/SDLW/development/pp/models"
	ae "github.com/vault-thirteen/auxie/errors"
)

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func readSettings(filePath string) (fps *models.FuncProcessorSettings, err error) {
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		derr := f.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	var fileContents []byte
	fileContents, err = io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	fps, err = models.NewFuncProcessorSettingsFromJson(fileContents)
	if err != nil {
		return nil, err
	}

	return fps, nil
}

func showFuncData(filePath string, fps *models.FuncProcessorSettings) (err error) {
	var fds []*models.FuncData
	fds, err = pp.GetFuncDataFromFile(filePath, fps)
	if err != nil {
		return err
	}

	fmt.Println("/* List of Go function names */")
	pp.ListGoFuncNames(fds)
	fmt.Println()
	fmt.Println("/* Manually added functions */")
	fmt.Println()
	fmt.Println("/* Automatically added functions */")
	fmt.Println()
	err = pp.PrintGoSysCallWrapperDrafts(fds)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	args, err := models.NewCliArgumentsFromOsArgs()
	mustBeNoError(err)

	var settings *models.FuncProcessorSettings
	settings, err = readSettings(args.FuncProcessorSettingsFilePath)
	mustBeNoError(err)

	err = showFuncData(args.HeaderFilePath, settings)
	mustBeNoError(err)
}
