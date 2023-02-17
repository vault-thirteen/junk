package user

import (
	"fmt"

	"github.com/vault-thirteen/auxie/unicode"
)

// Settings.
const (
	UserAuthenticationNameLengthMin     = 3
	UserAuthenticationNameLengthMax     = 40
	UserAuthenticationPasswordLengthMin = 8
	UserAuthenticationPasswordLengthMax = 64
)

// Errors.
const (
	ErrUserAuthenticationName     = "Error in UserAuthentication.Name: '%v'."
	ErrUserAuthenticationPassword = "Error in UserAuthentication.Password: '%v'."
)

// User Authentication Data.
type UserAuthentication struct {

	// Internal Name used for Authentication.
	Name string

	// Password used for Authentication.
	Password string
}

func ValidateUserAuthenticationName(
	name string,
) (err error) {
	if (len(name) < UserAuthenticationNameLengthMin) ||
		(len(name) > UserAuthenticationNameLengthMax) {
		err = fmt.Errorf(ErrUserAuthenticationName, name)
		return
	}
	var letters = []rune(name)
	if !unicode.SymbolIsLatLetter(letters[0]) {
		err = fmt.Errorf(ErrUserAuthenticationName, name)
		return
	}
	for i := 1; i < len(name); i++ {
		if (!unicode.SymbolIsLatLetter(letters[i])) &&
			(!unicode.SymbolIsNumber(letters[i])) {
			err = fmt.Errorf(ErrUserAuthenticationName, name)
			return
		}
	}
	return
}

func ValidateUserAuthenticationPassword(
	password string,
) (err error) {
	if (len(password) < UserAuthenticationPasswordLengthMin) ||
		(len(password) > UserAuthenticationPasswordLengthMax) {
		err = fmt.Errorf(ErrUserAuthenticationPassword, password)
		return
	}
	return
}
