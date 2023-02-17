package keysourcetype

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const KeySourceTypeNonExistent byte = 99

func TestIsValid(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedKeySourceType KeySourceType
		ResultExpected      bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1-A.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedKeySourceType: EnvironmentVariable,
		ResultExpected:      true,
	})

	// Тест 1-B.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedKeySourceType: File,
		ResultExpected:      true,
	})

	// Тест 1-C.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedKeySourceType: Vault,
		ResultExpected:      true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedKeySourceType: KeySourceType(KeySourceTypeNonExistent),
		ResultExpected:      false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedKeySourceType.IsValid()

		// Assert. Проверяем объект.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func TestIsEnvironmentVariable(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedKeySourceType KeySourceType
		ResultExpected      bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedKeySourceType: EnvironmentVariable,
		ResultExpected:      true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedKeySourceType: KeySourceType(KeySourceTypeNonExistent),
		ResultExpected:      false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedKeySourceType.IsEnvironmentVariable()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func TestIsFile(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedKeySourceType KeySourceType
		ResultExpected      bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedKeySourceType: File,
		ResultExpected:      true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedKeySourceType: KeySourceType(KeySourceTypeNonExistent),
		ResultExpected:      false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedKeySourceType.IsFile()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func TestIsVault(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedKeySourceType KeySourceType
		ResultExpected      bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedKeySourceType: Vault,
		ResultExpected:      true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedKeySourceType: KeySourceType(KeySourceTypeNonExistent),
		ResultExpected:      false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedKeySourceType.IsVault()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}

func Test_isEqualTo(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedKeySourceTypeFirst  KeySourceType
		TestedKeySourceTypeSecond KeySourceType
		ResultExpected            bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedKeySourceTypeFirst:  EnvironmentVariable,
		TestedKeySourceTypeSecond: EnvironmentVariable,
		ResultExpected:            true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedKeySourceTypeFirst:  EnvironmentVariable,
		TestedKeySourceTypeSecond: KeySourceType(KeySourceTypeNonExistent),
		ResultExpected:            false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedKeySourceTypeFirst.isEqualTo(test.TestedKeySourceTypeSecond)

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}
