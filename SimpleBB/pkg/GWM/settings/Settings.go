package s

import (
	"encoding/json"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"os"

	ver "github.com/vault-thirteen/auxie/Versioneer"
)

const (
	ErrUnknownClientIPAddressSource = "unknown client IP address source"
)

// Settings is Server's settings.
type settings struct {
	// Path to the file with these settings.
	FilePath simple.Path `json:"-"`

	// Program versioning information.
	VersionInfo *ver.Versioneer `json:"-"`

	IntHttpSettings  `json:"intHttp"`
	ExtHttpsSettings `json:"extHttps"`
	DbSettings       `json:"db"`
	ISystemSettings  `json:"system"`

	// External services.
	AcmSettings s.ServiceClientSettings `json:"acm"`
	MmSettings  s.ServiceClientSettings `json:"mm"`
	NmSettings  s.ServiceClientSettings `json:"nm"`
	SmSettings  s.ServiceClientSettings `json:"sm"`
}

func NewSettings() ISettings {
	return &settings{
		ISystemSettings: NewSystemSettings(),
	}
}

func NewSettingsFromFile(filePath string, versionInfo *ver.Versioneer) (stn ISettings, err error) {
	var buf []byte
	buf, err = os.ReadFile(filePath)
	if err != nil {
		return stn, err
	}

	stn = NewSettings()
	err = json.Unmarshal(buf, stn)
	if err != nil {
		return stn, err
	}

	stn.SetFilePath(simple.Path(filePath))

	err = stn.Check()
	if err != nil {
		return stn, err
	}

	dbs := stn.GetDbSettings()
	if len(dbs.Password) == 0 {
		var pwd string
		pwd, err = s.GetPasswordFromStdin(c.MsgEnterDatabasePassword)
		if err != nil {
			return stn, err
		}

		dbs.Password = pwd
		stn.SetDbSettings(dbs)
	}

	stn.SetVersionInfo(versionInfo)

	return stn, nil
}

func (stn *settings) Check() (err error) {
	err = s.CheckSettingsFilePath(stn.FilePath)
	if err != nil {
		return err
	}

	// Int. HTTP.
	err = stn.IntHttpSettings.Check()
	if err != nil {
		return err
	}

	// Ext. HTTPS.
	err = stn.ExtHttpsSettings.Check()
	if err != nil {
		return err
	}

	// DB.
	err = stn.DbSettings.Check()
	if err != nil {
		return err
	}

	// System.
	err = stn.GetSystemSettings().Check()
	if err != nil {
		return err
	}

	// External services.
	err = stn.AcmSettings.Check()
	if err != nil {
		return s.DetailedScsError(app.ServiceShortName_ACM, err)
	}

	err = stn.MmSettings.Check()
	if err != nil {
		return s.DetailedScsError(app.ServiceShortName_MM, err)
	}

	err = stn.NmSettings.Check()
	if err != nil {
		return s.DetailedScsError(app.ServiceShortName_NM, err)
	}

	err = stn.SmSettings.Check()
	if err != nil {
		return s.DetailedScsError(app.ServiceShortName_SM, err)
	}

	return nil
}

func (stn *settings) UseConstructor(filePath string, versionInfo *ver.Versioneer) (cmi.ISettings, error) {
	return NewSettingsFromFile(filePath, versionInfo)
}

// Emulated class members.
func (stn *settings) GetVersionInfo() (versionInfo *ver.Versioneer) {
	return stn.VersionInfo
}
func (stn *settings) GetDbSettings() (ds DbSettings) {
	return stn.DbSettings
}
func (stn *settings) GetIntHttpSettings() (ihs IntHttpSettings)   { return stn.IntHttpSettings }
func (stn *settings) GetExtHttpsSettings() (ehs ExtHttpsSettings) { return stn.ExtHttpsSettings }
func (stn *settings) GetSystemSettings() (ss ISystemSettings)     { return stn.ISystemSettings }
func (stn *settings) SetDbSettings(ds DbSettings) {
	stn.DbSettings = ds
}
func (stn *settings) SetFilePath(filePath simple.Path) {
	stn.FilePath = filePath
}
func (stn *settings) SetVersionInfo(versionInfo *ver.Versioneer) {
	stn.VersionInfo = versionInfo
}
func (s *settings) GetAcmSettings() s.ServiceClientSettings { return s.AcmSettings }
func (s *settings) GetMmSettings() s.ServiceClientSettings  { return s.MmSettings }
func (s *settings) GetNmSettings() s.ServiceClientSettings  { return s.NmSettings }
func (s *settings) GetSmSettings() s.ServiceClientSettings  { return s.SmSettings }
