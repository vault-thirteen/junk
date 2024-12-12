package u

import (
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	up "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/UserParameters"
)

// Unfortunately, Golang is a big lie. See comments below to know more.
type user struct {
	Password string `json:"-"`

	//TODO: Forget everything about Go language.
	//
	// Golang is not only shit, it is a big lie.
	// Here is the reason.
	//
	// https://stackoverflow.com/questions/38043678/golang-embedded-interface-on-parent-struct
	//
	// > Do not mix interface embedding with struct embedding.
	// > If you embed interfaces on a struct, you are actually adding new
	// > fields to the struct with the name of the interface so if you don't
	// > init those, you will get panics because they are nil.
	//
	// [Question 1]
	//
	// If embedding an interface is really not an embedding then why the hell
	// is it named as "embedding" ?
	//
	// [Question 2]
	//
	// Why embedding of structs works as real embedding, but "embedding" (fake
	// embedding) of interfaces does not work as embedding ?
	//
	// [Question 3]
	//
	// Is not it a lie when a fake embedding is called embedding ?

	base2.IUserParameters
}

func NewUser() (u derived1.IUser) {
	return &user{
		IUserParameters: up.NewUserParameters(),
	}
}

func NewUserFromScannableSource(src cmi.IScannable) (u derived1.IUser, err error) {
	u = NewUser()

	var uParams base2.IUserParameters
	uParams, err = up.NewUserParametersFromScannableSource(src)
	if err != nil {
		return nil, err
	}
	if uParams == nil {
		return nil, nil
	}

	//u.Password = "" // Password is never shown.
	u.SetUserParameters(uParams)

	return u, nil
}

// Emulated class members.
func (u *user) GetUserParameters() (up base2.IUserParameters) {
	return u.IUserParameters
}
func (u *user) SetUserParameters(up base2.IUserParameters) {
	u.IUserParameters = up
}
