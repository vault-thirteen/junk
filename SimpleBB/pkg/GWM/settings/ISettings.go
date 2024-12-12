package s

import (
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	s "github.com/vault-thirteen/SimpleBB/pkg/common/models/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	ver "github.com/vault-thirteen/auxie/Versioneer"
)

type ISettings interface {
	Check() (err error)
	UseConstructor(filePath string, versionInfo *ver.Versioneer) (cmi.ISettings, error)

	// Emulated class members.
	GetVersionInfo() (versionInfo *ver.Versioneer)
	GetDbSettings() (ds DbSettings)
	GetIntHttpSettings() (ihs IntHttpSettings)
	GetExtHttpsSettings() (ehs ExtHttpsSettings)
	GetSystemSettings() (ss ISystemSettings)
	SetDbSettings(ds DbSettings)
	SetFilePath(filePath simple.Path)
	SetVersionInfo(versionInfo *ver.Versioneer)
	GetAcmSettings() s.ServiceClientSettings
	GetMmSettings() s.ServiceClientSettings
	GetNmSettings() s.ServiceClientSettings
	GetSmSettings() s.ServiceClientSettings
}
