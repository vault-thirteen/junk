package application

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func setEnvVars() (err error) {
	err = os.Setenv(EnvVarNameIsDebugEnabled, "true")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameJwtKeySourceType, "1")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameJwtKeyDsn, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameJwtKeyValue, `-----BEGIN PUBLIC KEY-----.MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwjZtXrsyqk5zHxIG3toi.5drMsZW31rGaAVEY6E6dYsly1hNTb4kmk5J3IdhlwaJKqm0/0I0EVrVoZWwPKdgQ.xD8S3ekrMcSU4b6D6YInGOb5TrTLRkwlnBJzmfMZekngTadBb40hC+1ekQw2zln2.9e3Hmvn4hTKtY4AaG8dgiasd+ididnQQqhgZgdmJChkSvtoVcioPVGLGE9Yv6EbZ.7y/4aIWatMxrywOPoH85UHKT32XtAKtBzRLL/lvvBoeyzNCjZchQdm0fBcbC1yjI.Z3YCSdyMt+DfCKOy+BZYSHEpdX/dOyMt1rZFIroi0WpAt2xd6+z8W7rr+ru1QWrl.Tw0sBj3SM6mupsrtUOJsynY2sv6IIZ71huRmFvjEKWhrjb5A/s+XJp7dVd15noJv.UeCag3dJ+BM8Rd2TYoVy3F9sRBmhgh1v9l/eZYkTYiGBHMJu+gB3JtR3H/fVFACa.9vpkuOrm5Fy3p9xZsft0i80NtkY5Ad5e9tAA6ImhF+lp8Tkt+flJypqDnTNTGAH5.IVP283JGshzVDFOBi3xO25NiiznAXYmrfbhuxWErgnJZ+WbPCqTv7DaJ+v+d2vWx.EU7W9Yr4JH6Xtz7vRoyoZ35Nn6Fc3SWRDbvTCm675WacoEO3xkEqG4jMgxIlsWJf.Un51quvnqXEng6SDHvuXMNcCAwEAAQ==.-----END PUBLIC KEY-----.`)
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameKafkaConsumerGroupId, "gfe_test_consumer_group")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameKafkaBrokerAddressList, "127.0.0.1:9093")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameKafkaTopicList, "topic_a")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameStorageHost, "localhost")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStoragePort, "5432")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStorageUser, "test")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStoragePassword, "test")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStorageDatabase, "test")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStorageParameters, "")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameMetricsHttpServerHost, "0.0.0.0")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameMetricsHttpServerPort, "8888")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameBusinessLogicsHttpServerHost, "0.0.0.0")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameBusinessLogicsHttpServerPort, "9999")
	if err != nil {
		return err
	}

	return nil
}

func unsetEnvVars() (err error) {
	err = os.Setenv(EnvVarNameIsDebugEnabled, "")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameJwtKeySourceType, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameJwtKeyDsn, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameJwtKeyValue, "")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameKafkaConsumerGroupId, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameKafkaBrokerAddressList, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameKafkaTopicList, "")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameStorageHost, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStoragePort, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStorageUser, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStoragePassword, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStorageDatabase, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameStorageParameters, "")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameMetricsHttpServerHost, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameMetricsHttpServerPort, "")
	if err != nil {
		return err
	}

	err = os.Setenv(EnvVarNameBusinessLogicsHttpServerHost, "")
	if err != nil {
		return err
	}
	err = os.Setenv(EnvVarNameBusinessLogicsHttpServerPort, "")
	if err != nil {
		return err
	}

	return nil
}

func TestNewApplication(t *testing.T) {
	// Arrange.
	err := setEnvVars()
	assert.NoError(t, err)

	// Act.
	_, err = NewApplication()

	// Assert.
	assert.NoError(t, err)

	// Finalize.
	time.Sleep(time.Second)
	err = unsetEnvVars()
	assert.NoError(t, err)
}

func TestApplication_MustBeNoError(t *testing.T) {
	// Проверка этого метода невозможна.
}

func TestApplication_Start(t *testing.T) {
	// Arrange.
	err := setEnvVars()
	assert.NoError(t, err)
	app, err := NewApplication()
	assert.NoError(t, err)

	// Act.
	err = app.Start()

	// Assert.
	assert.NoError(t, err)

	// Finalize.
	time.Sleep(time.Second)
	err = app.Stop()
	assert.NoError(t, err)
	err = unsetEnvVars()
	assert.NoError(t, err)
}

func TestApplication_Stop(t *testing.T) {
	// Arrange.
	err := setEnvVars()
	assert.NoError(t, err)
	app, err := NewApplication()
	assert.NoError(t, err)
	err = app.Start()
	assert.NoError(t, err)
	time.Sleep(time.Second)

	// Act.
	err = app.Stop()

	// Assert.
	assert.NoError(t, err)

	// Finalize.
	err = unsetEnvVars()
	assert.NoError(t, err)
}

func TestApplication_WaitForQuitSignal(t *testing.T) {
	// Arrange.
	err := setEnvVars()
	assert.NoError(t, err)
	app, err := NewApplication()
	assert.NoError(t, err)
	err = app.Start()
	assert.NoError(t, err)
	time.Sleep(time.Second)
	app.quitSignals <- os.Interrupt

	// Act.
	err = app.WaitForQuitSignal()

	// Assert.
	assert.NoError(t, err)

	// Finalize.
	err = unsetEnvVars()
	assert.NoError(t, err)
}
