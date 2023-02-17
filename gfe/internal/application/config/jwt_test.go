package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vault-thirteen/junk/gfe/pkg/models/keysourcetype"
)

func TestNewJwt(t *testing.T) {
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
			Object: &Jwt{
				JwtKeySourceType: 0,
				JwtKeyDsn:        "",
				JwtKeyValue:      "",
			},
			IsErrorExpected: false,
		},
	})

	// Тест 2.
	// Переменные окружения заданы.
	testData = append(testData, TestDataForConfigConstructor{
		VariablesData: []TestEnvVarDataForConfigConstructor{
			{
				EnvVarName:  "TEST_JWT_KEY_SOURCE_TYPE",
				EnvVarValue: "123",
			},
			{
				EnvVarName:  "TEST_JWT_KEY_DSN",
				EnvVarValue: "dsn",
			},
			{
				EnvVarName:  "TEST_JWT_KEY_VALUE",
				EnvVarValue: "key",
			},
		},
		EnvPrefix: EnvPrefixCommon,
		ExpectedResult: ExpectedResultForConfigConstructor{
			Object: &Jwt{
				JwtKeySourceType: 123,
				JwtKeyDsn:        "dsn",
				JwtKeyValue:      "key",
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
		cfgActual, errActual := NewJwt(test.EnvPrefix)

		// Assert.1. Проверяем ошибку.
		switch test.ExpectedResult.IsErrorExpected {
		case true:
			assert.Error(t, errActual)
		case false:
			assert.NoError(t, errActual)
		}

		// Assert.2. Проверяем объект.
		expectedObject := test.ExpectedResult.Object.(*Jwt)
		assert.Equal(t, expectedObject, cfgActual)

		// Act.3. Убираем мусор из операционной системы.
		for _, envVar := range test.VariablesData {
			err := os.Setenv(envVar.EnvVarName, "")
			assert.NoError(t, err)
		}
	}
}

func TestJwt_IsValid(t *testing.T) {
	// Arrange.
	var testData = make([]TestDataForIsValid, 0)

	// Тест 1-A.
	// Объект с настройками полностью годен:
	// тип ключа = 1 (переменная окружения).
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Jwt{
			JwtKeySourceType: keysourcetype.EnvironmentVariable,
			JwtKeyDsn:        "",
			JwtKeyValue:      "key",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         true,
			IsErrorExpected: false,
		},
	})

	// Тест 1-B.
	// Объект с настройками полностью годен:
	// тип ключа = 2 (файл).
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Jwt{
			JwtKeySourceType: keysourcetype.File,
			JwtKeyDsn:        "dsn",
			JwtKeyValue:      "",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         true,
			IsErrorExpected: false,
		},
	})

	// Тест 1-C.
	// Объект с настройками полностью годен:
	// тип ключа = 3 (Vault).
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Jwt{
			JwtKeySourceType: keysourcetype.File,
			JwtKeyDsn:        "dsn",
			JwtKeyValue:      "",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         true,
			IsErrorExpected: false,
		},
	})

	// Тест 2-A-1.
	// Объект с настройками не годен:
	// тип ключа = 1 (переменная окружения), но DSN не пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Jwt{
			JwtKeySourceType: keysourcetype.EnvironmentVariable,
			JwtKeyDsn:        "junk",
			JwtKeyValue:      "key",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-A-2.
	// Объект с настройками не годен:
	// тип ключа = 1 (переменная окружения), но ключ пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Jwt{
			JwtKeySourceType: keysourcetype.EnvironmentVariable,
			JwtKeyDsn:        "",
			JwtKeyValue:      "",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-B-1.
	// Объект с настройками не годен:
	// тип ключа = 2 (файл), но DSN пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Jwt{
			JwtKeySourceType: keysourcetype.File,
			JwtKeyDsn:        "",
			JwtKeyValue:      "",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-B-2.
	// Объект с настройками не годен:
	// тип ключа = 2 (файл), но значение ключа указано.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Jwt{
			JwtKeySourceType: keysourcetype.File,
			JwtKeyDsn:        "dsn",
			JwtKeyValue:      "junk",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-C-1.
	// Объект с настройками не годен:
	// тип ключа = 3 (Vault), но но DSN пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Jwt{
			JwtKeySourceType: keysourcetype.Vault,
			JwtKeyDsn:        "",
			JwtKeyValue:      "",
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-C-2.
	// Объект с настройками не годен:
	// тип ключа = 3 (Vault), но но значение ключа указано.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &Jwt{
			JwtKeySourceType: keysourcetype.Vault,
			JwtKeyDsn:        "dsn",
			JwtKeyValue:      "junk",
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
		testedConfig := test.TestedConfig.(*Jwt)
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

func TestUnscreenJwtPublicKeyString(t *testing.T) {
	assert.Equal(t,
		"a\r\nb",
		UnScreenJwtPublicKeyString("a.b"),
	)
}

func TestScreenJwtPublicKeyString(t *testing.T) {
	assert.Equal(t,
		"a.b.c.d..e",
		ScreenJwtPublicKeyString("a\rb\nc\r\nd\n\re"),
	)
}
