package app

import (
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	ver "github.com/vault-thirteen/auxie/Versioneer"
)

func NewSettingsFromFile[T cmi.ISettings](classSelector T, filePath string, versionInfo *ver.Versioneer) (stn cmi.ISettings, err error) {
	return classSelector.UseConstructor(filePath, versionInfo)
}
