package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpServer(t *testing.T) {
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
			Object: &HttpServer{
				HttpServerHost: "0.0.0.0",
				HttpServerPort: 0,
			},
			IsErrorExpected: false,
		},
	})

	// Тест 2.
	// Переменные окружения заданы.
	testData = append(testData, TestDataForConfigConstructor{
		VariablesData: []TestEnvVarDataForConfigConstructor{
			{
				EnvVarName:  "TEST_HTTP_SERVER_HOST",
				EnvVarValue: "host",
			},
			{
				EnvVarName:  "TEST_HTTP_SERVER_PORT",
				EnvVarValue: "123",
			},
		},
		EnvPrefix: EnvPrefixCommon,
		ExpectedResult: ExpectedResultForConfigConstructor{
			Object: &HttpServer{
				HttpServerHost: "host",
				HttpServerPort: 123,
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
		cfgActual, errActual := NewHttpServer(test.EnvPrefix)

		// Assert.1. Проверяем ошибку.
		switch test.ExpectedResult.IsErrorExpected {
		case true:
			assert.Error(t, errActual)
		case false:
			assert.NoError(t, errActual)
		}

		// Assert.2. Проверяем объект.
		expectedObject := test.ExpectedResult.Object.(*HttpServer)
		assert.Equal(t, expectedObject, cfgActual)

		// Act.3. Убираем мусор из операционной системы.
		for _, envVar := range test.VariablesData {
			err := os.Setenv(envVar.EnvVarName, "")
			assert.NoError(t, err)
		}
	}
}

func TestHttpServer_IsValid(t *testing.T) {
	// Arrange.
	var testData = make([]TestDataForIsValid, 0)

	// Тест 1.
	// Объект с настройками полностью годен.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &HttpServer{
			HttpServerHost: "host",
			HttpServerPort: 12345,
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         true,
			IsErrorExpected: false,
		},
	})

	// Тест 2-A.
	// Объект с настройками не годен: кривой Хост.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &HttpServer{
			HttpServerHost: "",
			HttpServerPort: 12345,
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-B.
	// Объект с настройками не годен: кривой Порт.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &HttpServer{
			HttpServerHost: "host",
			HttpServerPort: 0,
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Запуск тестов.
	for i, test := range testData {
		fmt.Printf("[%d] ", i+1)

		// Act.1. Запускаем тестируемый метод или функцию.
		testedConfig := test.TestedConfig.(*HttpServer)
		isValidActual, errActual := testedConfig.IsValid()

		// Assert.1. Проверяем ошибку.
		switch test.ExpectedResult.IsErrorExpected {
		case true:
			assert.Error(t, errActual)
		case false:
			assert.NoError(t, errActual)
		}

		// Assert.2. Проверяем флаг.
		assert.Equal(t, test.ExpectedResult.IsValid, isValidActual)
	}
}
