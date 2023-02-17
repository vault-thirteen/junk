package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
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
			Object: &Storage{
				PostgreHost:       "localhost",
				PostgrePort:       5432,
				PostgreUser:       "",
				PostgrePassword:   "",
				PostgreDatabase:   "",
				PostgreParameters: "",
			},
			IsErrorExpected: false,
		},
	})

	// Тест 2.
	// Переменные окружения заданы.
	testData = append(testData, TestDataForConfigConstructor{
		VariablesData: []TestEnvVarDataForConfigConstructor{
			{
				EnvVarName:  "TEST_POSTGRE_HOST",
				EnvVarValue: "host",
			},
			{
				EnvVarName:  "TEST_POSTGRE_PORT",
				EnvVarValue: "12345",
			},
			{
				EnvVarName:  "TEST_POSTGRE_USER",
				EnvVarValue: "user",
			},
			{
				EnvVarName:  "TEST_POSTGRE_PASSWORD",
				EnvVarValue: "password",
			},
			{
				EnvVarName:  "TEST_POSTGRE_DATABASE",
				EnvVarValue: "database",
			},
			{
				EnvVarName:  "TEST_POSTGRE_PARAMETERS",
				EnvVarValue: "parameters",
			},
		},
		EnvPrefix: EnvPrefixCommon,
		ExpectedResult: ExpectedResultForConfigConstructor{
			Object: &Storage{
				PostgreHost:       "host",
				PostgrePort:       12345,
				PostgreUser:       "user",
				PostgrePassword:   "password",
				PostgreDatabase:   "database",
				PostgreParameters: "parameters",
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
		cfgActual, errActual := NewStorage(test.EnvPrefix)

		// Assert.1. Проверяем ошибку.
		switch test.ExpectedResult.IsErrorExpected {
		case true:
			assert.Error(t, errActual)
		case false:
			assert.NoError(t, errActual)
		}

		// Assert.2. Проверяем объект.
		expectedObject := test.ExpectedResult.Object.(*Storage)
		assert.Equal(t, expectedObject, cfgActual)

		// Act.3. Убираем мусор из операционной системы.
		for _, envVar := range test.VariablesData {
			err := os.Setenv(envVar.EnvVarName, "")
			assert.NoError(t, err)
		}
	}
}

func TestStorage_IsValid(t *testing.T) {
	// Arrange.
	var testData = make([]TestDataForIsValid, 0)

	// Тест 1.
	// Объект с настройками полностью годен.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Storage{
			PostgreHost:       "host",
			PostgrePort:       12345,
			PostgreUser:       "user",
			PostgrePassword:   "password",
			PostgreDatabase:   "database",
			PostgreParameters: "parameters",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         true,
			IsErrorExpected: false,
		},
	})

	// Тест 2-A.
	// Объект с настройками не годен:
	// хост пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Storage{
			PostgreHost:       "",
			PostgrePort:       12345,
			PostgreUser:       "user",
			PostgrePassword:   "password",
			PostgreDatabase:   "database",
			PostgreParameters: "parameters",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-B.
	// Объект с настройками не годен:
	// порт пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Storage{
			PostgreHost:       "host",
			PostgrePort:       0,
			PostgreUser:       "user",
			PostgrePassword:   "password",
			PostgreDatabase:   "database",
			PostgreParameters: "parameters",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-C.
	// Объект с настройками не годен:
	// пользователь пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Storage{
			PostgreHost:       "host",
			PostgrePort:       12345,
			PostgreUser:       "",
			PostgrePassword:   "password",
			PostgreDatabase:   "database",
			PostgreParameters: "parameters",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-D.
	// Объект с настройками не годен:
	// база данных не заполнена.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Storage{
			PostgreHost:       "host",
			PostgrePort:       12345,
			PostgreUser:       "user",
			PostgrePassword:   "password",
			PostgreDatabase:   "",
			PostgreParameters: "parameters",
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
		testedConfig := test.TestedConfig.(*Storage)
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
