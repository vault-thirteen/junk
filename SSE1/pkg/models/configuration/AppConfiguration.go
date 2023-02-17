package configuration

import (
	"errors"
	"time"

	vttime "github.com/vault-thirteen/auxie/time"
	configuration "github.com/vault-thirteen/junk/SSE1/pkg/models/configuration/xml"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/db/common"
)

// Errors.
const (
	ErrLoggerTypeUnknown  = "Logger Type is unknown"
	ErrStorageTypeUnknown = "Storage Type is unknown"
)

const (
	SourceFilePrefix = "file://"
)

type AppConfiguration struct {
	Source string
	Server ServerConfiguration
}

func NewAppConfiguration(
	filePath string,
) (appCfg *AppConfiguration, err error) {
	var xmlCfg *configuration.XmlAppConfiguration
	xmlCfg, err = configuration.NewXmlAppConfiguration(filePath)
	if err != nil {
		return
	}
	appCfg = &AppConfiguration{
		Source: SourceFilePrefix + filePath,
		Server: ServerConfiguration{
			Access: ServerAccessConfiguration{
				CoolDownPeriod: ServerAccessCoolDownPeriod{
					UserLogInSec: xmlCfg.Server.Access.CoolDownPeriod.UserLogInSec,
					UserUnregSec: xmlCfg.Server.Access.CoolDownPeriod.UserUnregSec,
				},
				Session: ServerAccessSessionConfiguration{
					IdleSessionTimeoutSec: xmlCfg.Server.Access.Session.IdleSessionTimeoutSec,
				},
				Token: ServerAccessTokenConfiguration{
					LifeTimeSec: xmlCfg.Server.Access.Token.LifeTimeSec,
				},
			},
			HttpServer: ServerHttpServerConfiguration{
				Address:            xmlCfg.Server.HttpServer.Address,
				CookiePath:         xmlCfg.Server.HttpServer.CookiePath,
				ShutdownTimeoutSec: xmlCfg.Server.HttpServer.ShutdownTimeoutSec,
				TokenHeader:        xmlCfg.Server.HttpServer.TokenHeader,
			},
			Logger: ServerLoggerConfiguration{
				IsEnabled: xmlCfg.Server.Logger.IsEnabled,
				Type:      NewServerLoggerType(xmlCfg.Server.Logger.Type),
			},
			Storage: ServerStorageConfiguration{
				Type: NewServerStorageType(xmlCfg.Server.Storage.Type),
				CommonParameters: common.StorageConfiguration{
					Address:              xmlCfg.Server.Storage.Address,
					ConnectionParameters: nil, // See below.
					Database:             xmlCfg.Server.Storage.Database,
					User:                 xmlCfg.Server.Storage.User,
					Password:             xmlCfg.Server.Storage.Password,
				},
				InitializationScripts: ServerStorageIniScripts{
					Folder: xmlCfg.Server.Storage.InitializationScripts.Folder,
				},
				//TableSettings: nil, // See below.
				Time: ServerStorageTimeConfiguration{
					//Zone: nil, // See below.
					Format: xmlCfg.Server.Storage.Time.Format,
				},
				// Settings taken from Application's Settings.
				CoolDownPeriods: ServerStorageCoolDownPeriods{
					UserLogIn: xmlCfg.Server.Access.CoolDownPeriod.UserLogInSec,
					UserUnreg: xmlCfg.Server.Access.CoolDownPeriod.UserUnregSec,
				},
				IdleSessionTimeoutSec: xmlCfg.Server.Access.Session.IdleSessionTimeoutSec,
				TokenLifeTimeSec:      xmlCfg.Server.Access.Token.LifeTimeSec,
			},
			Timezone: ServerTimezone{
				Name: xmlCfg.Server.TimeZone,
			},
			TLS: ServerTlsConfiguration{
				CertificateFile: xmlCfg.Server.TLS.CertificateFile,
				IsEnabled:       xmlCfg.Server.TLS.IsEnabled,
				KeyFile:         xmlCfg.Server.TLS.KeyFile,
			},
		},
	}
	if appCfg.Server.Logger.IsEnabled {
		if !appCfg.Server.Logger.Type.IsValid() {
			err = errors.New(ErrLoggerTypeUnknown)
			return
		}
	}
	appCfg.Server.Timezone.location, err = time.LoadLocation(
		appCfg.Server.Timezone.Name,
	)
	if err != nil {
		return
	}
	var hours int
	hours, err = vttime.GetLocationOffsetHours(
		appCfg.Server.Timezone.location,
	)
	if err != nil {
		return
	}
	appCfg.Server.Timezone.OffsetHrs = int8(hours)
	if !appCfg.Server.Storage.Type.IsValid() {
		err = errors.New(ErrStorageTypeUnknown)
		return
	}
	appCfg.Server.Storage.CommonParameters.ConnectionParameters, err =
		common.ParseConnectionParameters(xmlCfg.Server.Storage.ConnectionParameters)
	if err != nil {
		return
	}
	appCfg.Server.Storage.TableSettings, err = NewTableSettings(xmlCfg.Server.Storage.TableSettings)
	if err != nil {
		return
	}
	appCfg.Server.Storage.Time.Zone, err = time.LoadLocation(
		xmlCfg.Server.Storage.Time.Zone,
	)
	if err != nil {
		return
	}
	return
}
