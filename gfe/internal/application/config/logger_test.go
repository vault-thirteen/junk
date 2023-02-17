package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	// Внимание!
	// Нельзя задавать переменным окружения названия тех переменных, которые
	// уже используются в операционной системе.

	// Arrange.
	var testData = make([]TestDataForConfigConstructor, 0)

	// Тест 1.
	// Переменные окружения не заданы. Проверка значений по умолчанию.
	testData = append(testData, TestDataForConfigConstructor{
		VariablesData: []TestEnvVarDataForConfigConstructor{},
		EnvPrefix:     EnvPrefixCommon,
		ExpectedResult: ExpectedResultForConfigConstructor{
			Object: &Logger{
				IsDebugEnabled: false,
			},
			IsErrorExpected: false,
		},
	})

	// Тест 2.
	// Переменные окружения заданы.
	testData = append(testData, TestDataForConfigConstructor{
		VariablesData: []TestEnvVarDataForConfigConstructor{
			{
				EnvVarName:  "TEST_IS_DEBUG_ENABLED",
				EnvVarValue: "true",
			},
		},
		EnvPrefix: EnvPrefixCommon,
		ExpectedResult: ExpectedResultForConfigConstructor{
			Object: &Logger{
				IsDebugEnabled: true,
			},
			IsErrorExpected: false,
		},
	})

	// Запуск тестов.
	for i, test := range testData {
		fmt.Printf("[%d] ", i+1)

		// Act.1. Устанавливаем тестируемые переменные окружения.
		for _, envVar := range test.VariablesData {
			err := os.Setenv(envVar.EnvVarName, envVar.EnvVarValue)
			assert.NoError(t, err)
		}

		// Act.2. Запускаем тестируемый метод или функцию.
		cfgActual, errActual := NewLogger(test.EnvPrefix)

		// Assert.1. Проверяем ошибку.
		switch test.ExpectedResult.IsErrorExpected {
		case true:
			assert.Error(t, errActual)
		case false:
			assert.NoError(t, errActual)
		}

		// Assert.2. Проверяем объект.
		expectedObject := test.ExpectedResult.Object.(*Logger)
		assert.Equal(t, expectedObject, cfgActual)

		// Act.3. Убираем мусор из операционной системы.
		for _, envVar := range test.VariablesData {
			err := os.Setenv(envVar.EnvVarName, "")
			assert.NoError(t, err)
		}
	}
}
