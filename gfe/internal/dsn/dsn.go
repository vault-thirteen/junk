package dsn

import (
	"strings"

	"github.com/pkg/errors"
)

// Ошибки.
var (
	// ErrDsnIsNotFile -- ошибка "DSN -- не файл".
	ErrDsnIsNotFile = errors.Errorf("dsn is not a file")
)

// Форматы ошибок.
const (
	// ErrFDsnNotValid -- ошибка "DSN не годится".
	ErrFDsnNotValid = "dsn is not valid: %v"
)

// GetFilePathFromDsn извлекает путь файла из DSN (Data Source Name).
func GetFilePathFromDsn(dsn string) (filePath string, err error) {
	// schemaExpected -- ожидаемая схема.
	const schemaExpected = `file://`

	if !strings.HasPrefix(dsn, schemaExpected) {
		return "", ErrDsnIsNotFile
	}

	dsnParts := strings.SplitAfter(dsn, schemaExpected)
	if len(dsnParts) != 2 {
		return "", errors.Errorf(ErrFDsnNotValid, dsn)
	}

	filePath = dsnParts[1]

	if len(filePath) < 1 {
		return "", errors.Errorf(ErrFDsnNotValid, dsn)
	}

	return filePath, nil
}
