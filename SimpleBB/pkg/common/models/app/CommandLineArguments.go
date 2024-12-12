package app

import "strings"

type CommandLineArguments struct {
	configurationFilePathDefault string
	configurationFilePath        string
}

func NewCommandLineArgumentsFromOsArgs(osArgs []string, defaultCfgFilePath string) (cla *CommandLineArguments, err error) {
	cla = &CommandLineArguments{
		configurationFilePathDefault: defaultCfgFilePath,
	}

	if len(osArgs) != 2 {
		cla.configurationFilePath = cla.configurationFilePathDefault
		return cla, nil
	}

	cla.configurationFilePath = strings.TrimSpace(osArgs[1])

	return cla, nil
}

// IsDefaultFile tells whether the default file path is used for the
// configuration file.
func (cla *CommandLineArguments) IsDefaultFile() (isDefaultFile bool) {
	return cla.configurationFilePath == cla.configurationFilePathDefault
}

func (cla *CommandLineArguments) GetConfigurationFilePath() string {
	return cla.configurationFilePath
}
