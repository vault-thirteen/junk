package config

import "github.com/kelseyhightower/envconfig"

// Logger -- настройка ведения журнала.
type Logger struct {
	// Включатель режима отладки.
	IsDebugEnabled bool `split_words:"true" default:"false"`
}

// NewLogger -- создатель настроек ведения журнала.
func NewLogger(envPrefix string) (cfg *Logger, err error) {
	cfg = new(Logger)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsValid проверяет годность настройки ведения журнала.
func (l *Logger) IsValid() (bool, error) {
	return true, nil
}

// GetLoggerConfig получает настройки ведения журнала.
func GetLoggerConfig() (loggerConfig *Logger, err error) {
	loggerConfig, err = NewLogger(environmentVariablePrefixApplication)
	if err != nil {
		return nil, err
	}

	_, err = loggerConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return loggerConfig, nil
}
