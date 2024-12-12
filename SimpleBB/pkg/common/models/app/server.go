package app

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
)

func NewServer[T base.IServer](classSelector T, settings base.ISettings) (srv base.IServer, err error) {
	return classSelector.UseConstructor(settings)
}
