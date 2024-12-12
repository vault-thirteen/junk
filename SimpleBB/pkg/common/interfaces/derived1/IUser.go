package derived1

import "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"

type IUser interface {
	// Emulated class members.
	GetUserParameters() (up base2.IUserParameters)
	SetUserParameters(up base2.IUserParameters)
}
