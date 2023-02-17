package event

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const IDNonExistent = 999

func TestIsValid(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedID        TypeID
		ResultExpected  bool
		IsErrorExpected bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1-A.
	// Ожидается {true,nil}.
	tests = append(tests, TestData{
		TestedID:        Creation,
		ResultExpected:  true,
		IsErrorExpected: false,
	})

	// Тест 1-B.
	// Ожидается {true,nil}.
	tests = append(tests, TestData{
		TestedID:        Upload,
		ResultExpected:  true,
		IsErrorExpected: false,
	})

	// Тест 1-C.
	// Ожидается {true,nil}.
	tests = append(tests, TestData{
		TestedID:        Download,
		ResultExpected:  true,
		IsErrorExpected: false,
	})

	// Тест 1-D.
	// Ожидается {true,nil}.
	tests = append(tests, TestData{
		TestedID:        Modification,
		ResultExpected:  true,
		IsErrorExpected: false,
	})

	// Тест 2.
	// Ожидается {false,non-nil}.
	tests = append(tests, TestData{
		TestedID:        TypeID(IDNonExistent),
		ResultExpected:  false,
		IsErrorExpected: true,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual, errActual := test.TestedID.IsValid()

		// Assert.1. Проверяем ошибку.
		switch test.IsErrorExpected {
		case true:
			assert.Error(t, errActual)
		case false:
			assert.NoError(t, errActual)
		}

		// Assert.2. Проверяем объект.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func TestIsCreation(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedID       TypeID
		ResultExpected bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedID:       Creation,
		ResultExpected: true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedID:       TypeID(IDNonExistent),
		ResultExpected: false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedID.IsCreation()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func TestIsUpload(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedID       TypeID
		ResultExpected bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedID:       Upload,
		ResultExpected: true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedID:       TypeID(IDNonExistent),
		ResultExpected: false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedID.IsUpload()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func TestIsDownload(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedID       TypeID
		ResultExpected bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedID:       Download,
		ResultExpected: true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedID:       TypeID(IDNonExistent),
		ResultExpected: false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedID.IsDownload()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func TestIsModification(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedID       TypeID
		ResultExpected bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedID:       Modification,
		ResultExpected: true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedID:       TypeID(IDNonExistent),
		ResultExpected: false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedID.IsModification()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func TestIsSimple(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedID       TypeID
		ResultExpected bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1-A.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedID:       Creation,
		ResultExpected: true,
	})

	// Тест 1-B.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedID:       Upload,
		ResultExpected: true,
	})

	// Тест 1-C.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedID:       Modification,
		ResultExpected: true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedID:       TypeID(IDNonExistent),
		ResultExpected: false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedID.IsSimple()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func TestIsAggregated(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedID       TypeID
		ResultExpected bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedID:       Download,
		ResultExpected: true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedID:       TypeID(IDNonExistent),
		ResultExpected: false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedID.IsAggregated()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func Test_isEqualTo(t *testing.T) {
	const TestedID = TypeID(123)

	// TestData -- тестовые данные.
	type TestData struct {
		TestedIDFirst  TypeID
		TestedIDSecond TypeID
		ResultExpected bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedIDFirst:  TestedID,
		TestedIDSecond: TestedID,
		ResultExpected: true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedIDFirst:  TestedID,
		TestedIDSecond: IDNonExistent,
		ResultExpected: false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedIDFirst.isEqualTo(test.TestedIDSecond)

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}
