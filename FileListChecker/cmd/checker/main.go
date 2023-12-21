package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	stn "github.com/vault-thirteen/FileNameChecker/pkg/Settings"
	ver "github.com/vault-thirteen/auxie/Versioneer"
	ae "github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/file"
	ar "github.com/vault-thirteen/auxie/reader"
)

const (
	ErrfFileDoesNotExist = "file does not exist: %v"
)

const (
	MsgfAllFilesExist = "all files from the list in '%v' do exist."
)

func main() {
	showIntro()

	var settings stn.Settings
	var err error
	settings, err = getSettingsFromCommandLine()
	mustBeNoError(err)
	err = checkFiles(settings)
	mustBeNoError(err)
	var msg = fmt.Sprintf(MsgfAllFilesExist, settings.FileWithNamesPath)
	fmt.Print(msg)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func showIntro() {
	versioneer, err := ver.New()
	mustBeNoError(err)
	versioneer.ShowIntroText("")
	versioneer.ShowComponentsInfoText()
	fmt.Println()
}

func getSettingsFromCommandLine() (settings stn.Settings, err error) {
	var clArgs = os.Args
	if len(clArgs) < 2 {
		return settings, errors.New(stn.ErrCLAFolderNotSet)
	}
	if len(clArgs) < 3 {
		return settings, errors.New(stn.ErrCLAFileWithNamesNotSet)
	}

	settings.FolderPath = clArgs[1]
	settings.FileWithNamesPath = clArgs[2]

	return settings, nil
}

func checkFiles(settings stn.Settings) (err error) {
	var fileNames []string
	fileNames, err = getFileNames(settings.FileWithNamesPath)
	if err != nil {
		return err
	}

	var filePath string
	for _, fileName := range fileNames {
		filePath = filepath.Join(settings.FolderPath, fileName)
		err = ensureThatFileExists(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFileNames(fileWithNamesPath string) (fileNames []string, err error) {
	var f *os.File
	f, err = os.Open(fileWithNamesPath)
	if err != nil {
		return fileNames, err
	}

	defer func() {
		derr := f.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	var reader = ar.New(f)
	var bytes []byte
	bytes, err = reader.ReadLineEndingWithCRLF()
	for {
		if err != nil {
			if err == io.EOF {
				return fileNames, nil
			}
			return fileNames, err
		}

		filePath := strings.TrimSpace(string(bytes))
		if len(filePath) > 0 {
			fileNames = append(fileNames, filePath)
		}

		bytes, err = reader.ReadLineEndingWithCRLF()
	}
}

func ensureThatFileExists(filePath string) (err error) {
	var exists bool
	exists, err = file.Exists(filePath)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf(ErrfFileDoesNotExist, filePath)
	}

	return nil
}
