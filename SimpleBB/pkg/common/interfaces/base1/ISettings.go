package base

import (
	ver "github.com/vault-thirteen/auxie/Versioneer"
)

type ISettings interface {
	Check() error
	UseConstructor(string, *ver.Versioneer) (ISettings, error)
}
