package helper

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"spdb_manager/models"
	"spdb_manager/settings"
	"strconv"
	"strings"

	"go.uber.org/multierr"
)

// GetFilesList fetches the list of files using the HTTP GET method.
func GetFilesList(videosListUrl string) (paths []string, err error) {
	var resp *http.Response
	resp, err = http.Get(videosListUrl)
	if err != nil {
		return []string{}, err
	}

	defer func() {
		derr := resp.Body.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	var fileContents []byte
	fileContents, err = io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	paths = make([]string, 0)
	reader := strings.NewReader(string(fileContents))
	lineScanner := bufio.NewScanner(reader)

	for lineScanner.Scan() {
		paths = append(paths, lineScanner.Text())
	}

	return paths, nil
}

// GetFileTitle transforms a file name into a more human-friendly text.
func GetFileTitle(fileName string) (title string) {
	title = strings.TrimSuffix(fileName, path.Ext(fileName))
	title = strings.Replace(title, "_", " ", -1)
	title = strings.Replace(title, "-", " ", -1)
	title = strings.Replace(title, ".", " ", -1)
	for {
		if !strings.Contains(title, "  ") {
			break
		}
		title = strings.Replace(title, "  ", " ", -1)
	}

	return title
}

// GetDirectoryFilesList lists files in a folder.
func GetDirectoryFilesList(folder string) (files []string, err error) {
	var info []fs.FileInfo
	info, err = ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	files = make([]string, 0, len(info))
	for _, file := range info {
		if !file.IsDir() {
			files = append(files, file.Name())
		}
	}

	return files, nil
}

// GetFileContents returns file's contents as an array of bytes.
func GetFileContents(filePath string) (data []byte, err error) {
	var file *os.File
	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := file.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	data, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetFileContentsText returns file's contents as string.
func GetFileContentsText(filePath string) (data string, err error) {
	var tmp []byte
	tmp, err = GetFileContents(filePath)
	if err != nil {
		return "", err
	}

	return string(tmp), nil
}

// GetSqlLastInsertId returns the last inserted ID of an SQL operation.
func GetSqlLastInsertId(result sql.Result) (id models.ID, err error) {
	var lastInsertId int64
	lastInsertId, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return models.ID(lastInsertId), nil
}

// GetQueryParameter gets a string query parameter by its name.
// This method should be used only for a single parameter.
// If number of parameters is big, this method is not effective.
func GetQueryParameter(r *http.Request, paramName string) (paramValue string) {
	return r.URL.Query().Get(paramName)
}

// GetQueryParameterAsId gets an ID query parameter by its name.
// This method should be used only for a single parameter.
// If number of parameters is big, this method is not effective.
func GetQueryParameterAsId(r *http.Request, paramName string) (paramValue models.ID, err error) {
	valueString := GetQueryParameter(r, paramName)

	var valueUint64 uint64
	valueUint64, err = strconv.ParseUint(valueString, 10, 64)
	if err != nil {
		return 0, err
	}

	return models.ID(valueUint64), nil
}

// MakeTableNameForTagNMovies composes the name for a table containing movie
// indices for the tag.
func MakeTableNameForTagNMovies(tagId models.ID) (tableName string) {
	return fmt.Sprintf(settings.TableNameFormat_TagNMovies, tagId)
}
