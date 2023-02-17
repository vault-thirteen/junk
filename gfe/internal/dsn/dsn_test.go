package dsn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFilePathFromDsn(t *testing.T) {
	type TestData struct {
		Dsn              string
		FilePathExpected string
		IsErrorExpected  bool
	}

	// Arrange.
	var tests = make([]TestData, 0)

	// Тест 1-A.
	// DSN содержит ошибку -- кривой префикс.
	tests = append(tests, TestData{
		Dsn:              "junk",
		FilePathExpected: "",
		IsErrorExpected:  true,
	})

	// Тест 1-B.
	// DSN содержит ошибку -- количество частей неверно.
	tests = append(tests, TestData{
		Dsn:              "file://x file://y",
		FilePathExpected: "",
		IsErrorExpected:  true,
	})

	// Тест 1-C.
	// DSN содержит ошибку -- пустота.
	tests = append(tests, TestData{
		Dsn:              "file://",
		FilePathExpected: "",
		IsErrorExpected:  true,
	})

	// Тест 2.
	// Ошибок нет.
	tests = append(tests, TestData{
		Dsn:              "file://file.txt",
		FilePathExpected: "file.txt",
		IsErrorExpected:  false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		filePathActual, errActual := GetFilePathFromDsn(test.Dsn)

		// Assert.1. Проверяем ошибку.
		switch test.IsErrorExpected {
		case true:
			assert.Error(t, errActual)
		case false:
			assert.NoError(t, errActual)
		}

		// Assert.2. Проверяем объект.
		assert.Equal(t, test.FilePathExpected, filePathActual)
	}
}
