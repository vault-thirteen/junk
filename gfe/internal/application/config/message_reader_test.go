package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageReader(t *testing.T) {
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
			Object: &MessageReader{
				KafkaConsumerGroupID:   "",
				KafkaBrokerAddressList: nil,
				KafkaTopicList:         nil,
			},
			IsErrorExpected: false,
		},
	})

	// Тест 2.
	// Переменные окружения заданы.
	testData = append(testData, TestDataForConfigConstructor{
		VariablesData: []TestEnvVarDataForConfigConstructor{
			{
				EnvVarName:  "TEST_KAFKA_CONSUMER_GROUP_ID",
				EnvVarValue: "id",
			},
			{
				EnvVarName:  "TEST_KAFKA_BROKER_ADDRESS_LIST",
				EnvVarValue: "address_1,address_2,address_3",
			},
			{
				EnvVarName:  "TEST_KAFKA_TOPIC_LIST",
				EnvVarValue: "topic_1,topic_2,topic_3",
			},
		},
		EnvPrefix: EnvPrefixCommon,
		ExpectedResult: ExpectedResultForConfigConstructor{
			Object: &MessageReader{
				KafkaConsumerGroupID: "id",
				KafkaBrokerAddressList: []string{
					"address_1",
					"address_2",
					"address_3",
				},
				KafkaTopicList: []string{
					"topic_1",
					"topic_2",
					"topic_3",
				},
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
		cfgActual, errActual := NewMessageReader(test.EnvPrefix)

		// Assert.1. Проверяем ошибку.
		switch test.ExpectedResult.IsErrorExpected {
		case true:
			assert.Error(t, errActual)
		case false:
			assert.NoError(t, errActual)
		}

		// Assert.2. Проверяем объект.
		expectedObject := test.ExpectedResult.Object.(*MessageReader)
		assert.Equal(t, expectedObject, cfgActual)

		// Act.3. Убираем мусор из операционной системы.
		for _, envVar := range test.VariablesData {
			err := os.Setenv(envVar.EnvVarName, "")
			assert.NoError(t, err)
		}
	}
}

func TestMessageReader_IsValid(t *testing.T) {
	// Arrange.
	var testData = make([]TestDataForIsValid, 0)

	// Тест 1.
	// Объект с настройками полностью годен.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &MessageReader{
			KafkaConsumerGroupID:   "id",
			KafkaBrokerAddressList: []string{"address_1"},
			KafkaTopicList:         []string{"topic_1"},
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         true,
			IsErrorExpected: false,
		},
	})

	// Тест 2-A.
	// Объект с настройками не годен:
	// ID группы потребителя Kafka пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &MessageReader{
			KafkaConsumerGroupID:   "",
			KafkaBrokerAddressList: []string{"address_1"},
			KafkaTopicList:         []string{"topic_1"},
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-B.
	// Объект с настройками не годен:
	// список адресов посредников (брокеров) Kafka пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &MessageReader{
			KafkaConsumerGroupID:   "id",
			KafkaBrokerAddressList: []string{},
			KafkaTopicList:         []string{"topic_1"},
		},
		ExpectedResult: ExpectedResultForIsValid{
			IsValid:         false,
			IsErrorExpected: true,
		},
	})

	// Тест 2-C.
	// Объект с настройками не годен:
	// список тем (топиков) Kafka пуст.
	testData = append(testData, TestDataForIsValid{
		TestedConfig: &MessageReader{
			KafkaConsumerGroupID:   "id",
			KafkaBrokerAddressList: []string{"address_1"},
			KafkaTopicList:         []string{},
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
		testedConfig := test.TestedConfig.(*MessageReader)
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
