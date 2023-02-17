package file

import (
	"io"
	"os"

	"github.com/vault-thirteen/errorz"
)

// Gets Contents of the Text File.
// This Method is more safe than the built-in 'ioutil.ReadFile' which may skip
// some Errors.
func GetTextFileContents(
	filePath string,
) (contents string, err error) {
	var file *os.File
	file, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer func() {
		var derr = file.Close()
		if derr != nil {
			err = errorz.Combine(err, derr)
		}
	}()

	var contentsBA []byte
	contentsBA, err = io.ReadAll(file)
	if err != nil {
		return
	}
	contents = string(contentsBA)
	return
}
