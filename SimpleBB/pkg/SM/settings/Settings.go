package s

import (
	"encoding/json"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/settings"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"os"

	ver "github.com/vault-thirteen/auxie/Versioneer"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the file with these settings.
	FilePath cm.Path `json:"-"`

	// Program versioning information.
	VersionInfo *ver.Versioneer `json:"-"`

	HttpsSettings  `json:"https"`
	DbSettings     `json:"db"`
	SystemSettings `json:"system"`

	// External services.
	AcmSettings s.ServiceClientSettings `json:"acm"`
	MmSettings  s.ServiceClientSettings `json:"mm"`
}

func NewSettingsFromFile(filePath string, versionInfo *ver.Versioneer) (stn *Settings, err error) {
	var buf []byte
	buf, err = os.ReadFile(filePath)
	if err != nil {
		return stn, err
	}

	stn = &Settings{}
	err = json.Unmarshal(buf, stn)
	if err != nil {
		return stn, err
	}

	stn.FilePath = cm.Path(filePath)

	err = stn.Check()
	if err != nil {
		return stn, err
	}

	if len(stn.Password) == 0 {
		stn.DbSettings.Password, err = s.GetPasswordFromStdin(c.MsgEnterDatabasePassword)
		if err != nil {
			return stn, err
		}
	}

	stn.VersionInfo = versionInfo

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	err = s.CheckSettingsFilePath(stn.FilePath)
	if err != nil {
		return err
	}

	// HTTPS.
	err = stn.HttpsSettings.Check()
	if err != nil {
		return err
	}

	// DB.
	err = stn.DbSettings.Check()
	if err != nil {
		return err
	}

	// System.
	err = stn.SystemSettings.Check()
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

	return nil
}

func (stn *Settings) UseConstructor(filePath string, versionInfo *ver.Versioneer) (cmi.ISettings, error) {
	return NewSettingsFromFile(filePath, versionInfo)
}
