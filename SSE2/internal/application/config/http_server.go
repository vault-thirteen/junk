package config

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	HttpServerShutdownTimeoutSec = 60
)

type HttpServer struct {
	HttpServerHost string `split_words:"true" default:"0.0.0.0"`
	HttpServerPort uint   `split_words:"true"`
}

func NewHttpServer(envPrefix string) (cfg *HttpServer, err error) {
	cfg = new(HttpServer)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (hsc *HttpServer) IsValid() (bool, error) {
	if len(hsc.HttpServerHost) < 1 {
		return false, ErrHost
	}

	if hsc.HttpServerPort < 1 {
		return false, ErrPort
	}

	return true, nil
}

func GetHttpServerConfig() (httpServerConfig *HttpServer, err error) {
	httpServerConfig, err = NewHttpServer(EnvironmentVariablePrefixApplication)
	if err != nil {
		return nil, err
	}

	_, err = httpServerConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return httpServerConfig, nil
}
