package cm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"

	"github.com/vault-thirteen/TR1/src/interfaces"
)

const (
	OsArgsCountExpected          = 1
	ConfigurationFilePathDefault = "settings.json"
)

const (
	Err_ArgsCount                 = "number of command line arguments is incorrect"
	Err_ServiceShortNameIsNotSet  = "service short name is not set"
	ErrF_ServiceShortNameMismatch = "service short name mismatch: %s vs %s"
)

const (
	Msg_ServiceIsStarting     = "Service is starting ..."
	Msg_ServiceIsStopping     = "Service is stopping ..."
	Msg_ServiceIsStopped      = "Service is stopped"
	MsgF_QuitSignalIsReceived = "Quit signal from OS has been received: %v"
)

type Application struct {
	shortName             string
	configurationFilePath string
	ver                   *ver.Versioneer
	cfg                   *Configuration
	service               *Service
	dl                    *DllLoader
}

func NewApplication(serviceShortName string, serviceComponents []interfaces.IServiceComponent, controller interfaces.IController) (a *Application, err error) {
	if len(serviceShortName) == 0 {
		return nil, errors.New(Err_ServiceShortNameIsNotSet)
	}

	var configurationFilePath string
	configurationFilePath, err = getConfigurationFilePathFromOs()
	if err != nil {
		return nil, err
	}

	var vi *ver.Versioneer
	vi, err = ver.New(false)
	if err != nil {
		return nil, err
	}

	var cfg *Configuration
	cfg, err = NewConfigurationFromFile(configurationFilePath)
	if err != nil {
		return nil, err
	}

	if cfg.Service.ShortName != serviceShortName {
		return nil, fmt.Errorf(ErrF_ServiceShortNameMismatch, cfg.Service.ShortName, serviceShortName)
	}

	a = &Application{
		shortName:             serviceShortName,
		configurationFilePath: configurationFilePath,
		ver:                   vi,
		cfg:                   cfg,
		dl:                    NewDllLoader(),
	}

	showIntro(a.ver, a.shortName)

	err = a.dl.Init()
	if err != nil {
		return nil, err
	}

	a.service, err = NewService(a.cfg, serviceComponents, controller)
	if err != nil {
		return nil, err
	}

	err = controller.LinkWithService(a.service)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func getConfigurationFilePathFromOs() (cfp string, err error) {
	switch len(os.Args) {
	case 1:
		// No arguments are set.
		// Use the default path for configuration file.
		return ConfigurationFilePathDefault, nil

	case OsArgsCountExpected + 1:
		cfp = strings.TrimSpace(os.Args[1])
		return cfp, nil

	default:
		return "", errors.New(Err_ArgsCount)
	}
}

func (a Application) GetConfiguration() interfaces.IConfiguration {
	return *a.cfg
}

func (a *Application) Use() (err error) {
	// Start.
	log.Println(Msg_ServiceIsStarting)

	err = a.service.Start()
	if err != nil {
		return err
	}

	a.service.ReportStart()

	// Run.
	serverMustBeStopped := a.service.GetStopChannel()
	waitForQuitSignalFromOS(serverMustBeStopped)
	<-*serverMustBeStopped

	// Stop.
	log.Println(Msg_ServiceIsStopping)

	err = a.service.Stop()
	if err != nil {
		return err
	}

	log.Println(Msg_ServiceIsStopped)
	time.Sleep(time.Second)

	return nil
}

func showIntro(v *ver.Versioneer, serviceName string) {
	v.ShowIntroText(serviceName)
	v.ShowComponentsInfoText()
	fmt.Println()
}

func waitForQuitSignalFromOS(serverMustBeStopped *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range osSignals {
			switch sig {
			case syscall.SIGINT,
				syscall.SIGTERM:
				log.Println(fmt.Sprintf(MsgF_QuitSignalIsReceived, sig))
				*serverMustBeStopped <- true
			}
		}
	}()
}
