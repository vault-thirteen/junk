package fileeventsrequest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRecordsCountLimitSet(t *testing.T) {
	// TestData -- тестовые данные.
	type TestData struct {
		TestedObject   *FileEventsRequest
		ResultExpected bool
	}

	// Arrange.
	tests := make([]TestData, 0)

	// Тест 1.
	// Ожидается true.
	tests = append(tests, TestData{
		TestedObject: &FileEventsRequest{
			RecordsCountLimit: 123,
		},
		ResultExpected: true,
	})

	// Тест 2.
	// Ожидается false.
	tests = append(tests, TestData{
		TestedObject: &FileEventsRequest{
			RecordsCountLimit: 0,
		},
		ResultExpected: false,
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		resultActual := test.TestedObject.IsRecordsCountLimitSet()

		// Assert. Проверяем результат.
		assert.Equal(t, test.ResultExpected, resultActual)
	}
}
