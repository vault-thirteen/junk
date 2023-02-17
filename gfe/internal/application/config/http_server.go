package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/vault-thirteen/junk/gfe/internal/helper"
)

// Различные настройки.
const (
	// HttpServerShutdownTimeoutSec -- предельный период времени (в секундах),
	// в течение которого HTTP сервер должен выключаться.
	HttpServerShutdownTimeoutSec = 60
)

// HttpServer -- настройка HTTP сервера.
type HttpServer struct {
	// Прослушиваемый хост.
	HttpServerHost string `split_words:"true" default:"0.0.0.0"`

	// Прослушиваемый порт.
	HttpServerPort uint `split_words:"true"`
}

// NewHttpServer -- создатель настроек HTTP сервера.
func NewHttpServer(envPrefix string) (cfg *HttpServer, err error) {
	cfg = new(HttpServer)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsValid проверяет годность настройки HTTP сервера.
func (hsc *HttpServer) IsValid() (bool, error) {
	if len(hsc.HttpServerHost) < 1 {
		return false, ErrHost
	}

	if hsc.HttpServerPort < 1 {
		return false, ErrPort
	}

	return true, nil
}

// GetBusinessHttpServerConfig получает настройки HTTP сервера бизнес логики.
func GetBusinessHttpServerConfig() (httpServerConfig *HttpServer, err error) {
	envPrefix := helper.ConcatenateEnvVarPrefixes(
		environmentVariablePrefixApplication,
		environmentVariablePrefixBusinessLogicsServer,
	)

	httpServerConfig, err = NewHttpServer(envPrefix)
	if err != nil {
		return nil, err
	}

	_, err = httpServerConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return httpServerConfig, nil
}

// GetSystemHttpServerConfig получает настройки системного HTTP сервера.
func GetSystemHttpServerConfig() (httpServerConfig *HttpServer, err error) {
	envPrefix := helper.ConcatenateEnvVarPrefixes(
		environmentVariablePrefixApplication,
		environmentVariablePrefixMetricsServer,
	)

	httpServerConfig, err = NewHttpServer(envPrefix)
	if err != nil {
		return nil, err
	}

	_, err = httpServerConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return httpServerConfig, nil
}
