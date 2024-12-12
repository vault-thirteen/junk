package app

import (
	"errors"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"log"
	"os"
	"time"

	ver "github.com/vault-thirteen/auxie/Versioneer"
)

const (
	ErrServiceName                  = "service name error"
	ErrConfigurationFilePathDefault = "default configuration file path error"
)

type Application struct {
	serviceName                  string
	configurationFilePathDefault string
	ver                          *ver.Versioneer
	cla                          *CommandLineArguments
	settings                     base.ISettings
	server                       base.IServer
}

func NewApplication[T1 base.ISettings, T2 base.IServer](
	settingsClassSelector T1,
	serverClassSelector T2,
	serviceName string,
	configurationFilePathDefault string,
) (a *Application, err error) {
	if len(serviceName) == 0 {
		return nil, errors.New(ErrServiceName)
	}

	if len(configurationFilePathDefault) == 0 {
		return nil, errors.New(ErrConfigurationFilePathDefault)
	}

	a = &Application{
		serviceName:                  serviceName,
		configurationFilePathDefault: configurationFilePathDefault,
	}

	a.ver, err = ver.New()
	if err != nil {
		return nil, err
	}

	a.cla, err = NewCommandLineArgumentsFromOsArgs(os.Args, a.configurationFilePathDefault)
	if err != nil {
		return nil, err
	}

	if a.cla.IsDefaultFile() {
		log.Println(c.MsgUsingDefaultConfigurationFile)
	}

	a.settings, err = NewSettingsFromFile[T1](settingsClassSelector, a.cla.GetConfigurationFilePath(), a.ver)
	if err != nil {
		return nil, err
	}

	a.server, err = NewServer[T2](serverClassSelector, a.settings)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) Use() (err error) {
	showIntro(a.ver, a.serviceName)

	// Start.
	log.Println(c.MsgServerIsStarting)

	err = a.server.Start()
	if err != nil {
		return err
	}

	a.server.ReportStart()

	// Run.
	serverMustBeStopped := a.server.GetStopChannel()
	waitForQuitSignalFromOS(serverMustBeStopped)
	<-*serverMustBeStopped

	// Stop.
	log.Println(c.MsgServerIsStopping)

	err = a.server.Stop()
	if err != nil {
		return err
	}

	log.Println(c.MsgServerIsStopped)
	time.Sleep(time.Second)

	return nil
}
