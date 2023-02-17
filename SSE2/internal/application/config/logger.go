package config

import "github.com/kelseyhightower/envconfig"

type Logger struct {
	IsDebugEnabled bool `split_words:"true" default:"false"`
}

func NewLogger(envPrefix string) (cfg *Logger, err error) {
	cfg = new(Logger)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (l *Logger) IsValid() (bool, error) {
	return true, nil
}

func GetLoggerConfig() (loggerConfig *Logger, err error) {
	loggerConfig, err = NewLogger(EnvironmentVariablePrefixApplication)
	if err != nil {
		return nil, err
	}

	_, err = loggerConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return loggerConfig, nil
}
