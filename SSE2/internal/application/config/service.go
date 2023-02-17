package config

import "github.com/kelseyhightower/envconfig"

type Service struct {
	WorkersCount                      uint   `split_words:"true"`
	PathToConverterExecutable         string `split_words:"true"`
	LargePngImageMaximumSideDimension uint   `split_words:"true"`
	SmallPngImageMaximumSideDimension uint   `split_words:"true"`
	FileSizeLimitSettingsFile         string `split_words:"true"`

	// https://bugs.documentfoundation.org/show_bug.cgi?id=37531.
	UseLibreOfficeMultipleUserInstallations bool `split_words:"true" default:"true"`
}

const (
	LibreOfficeTemporaryFolder        = "libreoffice"
	LibreOfficeUserInstallationFolder = "user_installation"
)

func NewService(envPrefix string) (cfg *Service, err error) {
	cfg = new(Service)
	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (s *Service) IsValid() (bool, error) {
	if s.WorkersCount < 1 {
		return false, ErrWorkersCount
	}

	if len(s.PathToConverterExecutable) < 1 {
		return false, ErrPathToConverterExecutable
	}

	if s.LargePngImageMaximumSideDimension < 1 {
		return false, ErrLargePngImageMaximumSideDimension
	}

	if s.SmallPngImageMaximumSideDimension < 1 {
		return false, ErrSmallPngImageMaximumSideDimension
	}

	if len(s.FileSizeLimitSettingsFile) < 1 {
		return false, ErrFileSizeLimitSettingsFile
	}

	return true, nil
}

func GetServiceConfig() (serviceConfig *Service, err error) {
	serviceConfig, err = NewService(EnvironmentVariablePrefixApplication)
	if err != nil {
		return nil, err
	}

	_, err = serviceConfig.IsValid()
	if err != nil {
		return nil, err
	}

	return serviceConfig, nil
}
