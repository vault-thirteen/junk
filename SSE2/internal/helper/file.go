package helper

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vault-thirteen/junk/SSE2/internal/random"
	"go.uber.org/multierr"
)

func MakeTemporaryFolderName() (folderName string) {
	return random.MakeUniqueRandomString()
}

func CreateSubFolder(
	parentFolderPath string,
	newFolderName string,
) (newFolderPath string, err error) {
	newFolderPath = filepath.Join(parentFolderPath, newFolderName)

	err = os.Mkdir(newFolderPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return newFolderPath, nil
}

func GetTextFileContents(
	filePath string,
) (fileContents string, err error) {
	var file *os.File
	file, err = os.Open(filePath)
	if err != nil {
		return
	}

	defer func() {
		var derr = file.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	var buffer []byte
	buffer, err = io.ReadAll(file)
	if err != nil {
		return
	}

	return string(buffer), nil
}

func MakeTemporaryLocalFileName(fileName string) (localFileName string) {
	return strings.ReplaceAll(fileName, `/`, `_`)
}
