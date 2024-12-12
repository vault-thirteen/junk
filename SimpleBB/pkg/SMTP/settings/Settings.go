package settings

import (
	"encoding/json"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/models/settings"
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

	HttpSettings   `json:"http"`
	SystemSettings `json:"system"`
	SmtpSettings   `json:"smtp"`
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
		stn.SmtpSettings.Password, err = cs.GetPasswordFromStdin(c.MsgEnterSmtpPassword)
		if err != nil {
			return stn, err
		}
	}

	stn.VersionInfo = versionInfo

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	err = cs.CheckSettingsFilePath(stn.FilePath)
	if err != nil {
		return err
	}

	// HTTP.
	err = stn.HttpSettings.Check()
	if err != nil {
		return err
	}

	// System.
	err = stn.SystemSettings.Check()
	if err != nil {
		return err
	}

	// SMTP.
	err = stn.SmtpSettings.Check()
	if err != nil {
		return err
	}

	return nil
}

func (stn *Settings) UseConstructor(filePath string, versionInfo *ver.Versioneer) (cmi.ISettings, error) {
	return NewSettingsFromFile(filePath, versionInfo)
}
