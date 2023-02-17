package integration

import (
	"os"

	"github.com/vault-thirteen/junk/gfe/internal/application/config"
)

// Список имён переменных окружения.
// Подробности можно увидеть в пакете 'config'.
const (
	EnvVarNameIsDebugEnabled = "GFE_IS_DEBUG_ENABLED"

	EnvVarNameJwtKeySourceType = "GFE_JWT_KEY_SOURCE_TYPE"
	EnvVarNameJwtKeyDsn        = "GFE_JWT_KEY_DSN"
	EnvVarNameJwtKeyValue      = "GFE_JWT_KEY_VALUE"

	EnvVarNameKafkaConsumerGroupId   = "GFE_KAFKA_CONSUMER_GROUP_ID"
	EnvVarNameKafkaBrokerAddressList = "GFE_KAFKA_BROKER_ADDRESS_LIST"
	EnvVarNameKafkaTopicList         = "GFE_KAFKA_TOPIC_LIST"

	EnvVarNameStorageHost       = "GFE_POSTGRE_HOST"
	EnvVarNameStoragePort       = "GFE_POSTGRE_PORT"
	EnvVarNameStorageUser       = "GFE_POSTGRE_USER"
	EnvVarNameStoragePassword   = "GFE_POSTGRE_PASSWORD"
	EnvVarNameStorageDatabase   = "GFE_POSTGRE_DATABASE"
	EnvVarNameStorageParameters = "GFE_POSTGRE_PARAMETERS"

	EnvVarNameMetricsHttpServerHost = "GFE_METRICS_HTTP_SERVER_HOST"
	EnvVarNameMetricsHttpServerPort = "GFE_METRICS_HTTP_SERVER_PORT"

	EnvVarNameBusinessLogicsHttpServerHost = "GFE_BUSINESS_HTTP_SERVER_HOST"
	EnvVarNameBusinessLogicsHttpServerPort = "GFE_BUSINESS_HTTP_SERVER_PORT"
)

// setEnvVars устанавливает переменные окружения из списка.
func (t *Test) setEnvVars(environmentVariables []EnvironmentVariable) (err error) {
	for _, environmentVariable := range environmentVariables {
		err = os.Setenv(environmentVariable.Name, environmentVariable.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// setEnvVars чистит переменные окружения из списка (устанавливает им пустое
// значение).
func (t *Test) unsetEnvVars() (err error) {
	for _, environmentVariable := range t.EnvironmentVariables {
		err = os.Setenv(environmentVariable.Name, "")
		if err != nil {
			return err
		}
	}

	return nil
}

// prepareEnvironmentVariablesList создаёт список имён и значений переменных
// окружения.
func (t *Test) prepareEnvironmentVariablesList() (vars []EnvironmentVariable, err error) {
	// Берём RSA ключ для JWT из заранее подготовленного файла и предаём его
	// через переменную окружения.
	var keyPemTextScreened string
	keyPemTextScreened, err = t.prepareRsaPublicKeyScreened(t.RsaPublicKeyFilePath)
	if err != nil {
		return nil, err
	}

	vars = []EnvironmentVariable{
		{
			Name:  EnvVarNameIsDebugEnabled,
			Value: "true",
		},
		{
			Name:  EnvVarNameJwtKeySourceType,
			Value: "1",
		},
		{
			Name:  EnvVarNameJwtKeyDsn,
			Value: "",
		},
		{
			Name:  EnvVarNameJwtKeyValue,
			Value: keyPemTextScreened,
		},
		{
			Name:  EnvVarNameKafkaConsumerGroupId,
			Value: composeId(IntegrationTestIdPrefix, "consumer_group"),
		},
		{
			Name:  EnvVarNameKafkaBrokerAddressList,
			Value: "127.0.0.1:9093",
		},
		{
			Name:  EnvVarNameKafkaTopicList,
			Value: composeTopicName(t.UniqueTestId),
		},
		{
			Name:  EnvVarNameStorageHost,
			Value: "localhost",
		},
		{
			Name:  EnvVarNameStoragePort,
			Value: "5432",
		},
		{
			Name:  EnvVarNameStorageUser,
			Value: IntegrationTestIdPrefix,
		},
		{
			Name:  EnvVarNameStoragePassword,
			Value: IntegrationTestIdPrefix,
		},
		{
			Name:  EnvVarNameStorageDatabase,
			Value: IntegrationTestIdPrefix,
		},
		{
			Name:  EnvVarNameStorageParameters,
			Value: "",
		},
		{
			Name:  EnvVarNameMetricsHttpServerHost,
			Value: "0.0.0.0",
		},
		{
			Name:  EnvVarNameMetricsHttpServerPort,
			Value: "2001",
		},
		{
			Name:  EnvVarNameBusinessLogicsHttpServerHost,
			Value: "0.0.0.0",
		},
		{
			Name:  EnvVarNameBusinessLogicsHttpServerPort,
			Value: "2002",
		},
	}

	return vars, nil
}

// prepareRsaPublicKeyScreened читает из файла публичный RSA ключ и экранирует
// его для использования в переменной окружения.
func (t *Test) prepareRsaPublicKeyScreened(keyFilePath string) (keyPemTextScreened string, err error) {
	var keyPemText string
	keyPemText, err = t.readRsaPublicKeyFromFile(keyFilePath)
	if err != nil {
		return "", err
	}

	return config.ScreenJwtPublicKeyString(keyPemText), nil
}

// readRsaPublicKeyFromFile читает из файла публичный RSA ключ.
func (t *Test) readRsaPublicKeyFromFile(filePath string) (keyPemText string, err error) {
	var buffer []byte
	buffer, err = os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}
