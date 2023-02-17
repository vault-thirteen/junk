package user

import "fmt"

// Settings.
const (
	UserPublicNameLengthMin = 3
	UserPublicNameLengthMax = 40
)

// Errors.
const (
	ErrUserPublicName = "Error in User.PublicName: '%v'."
)

// User.
type User struct {

	// 'true' means that a User is registered.
	// 'false' means that a User is disabled.
	IsEnabled bool

	// A unique Identifier in the List of all Users.
	Id uint

	// Authentication Data.
	Authentication UserAuthentication

	// Registration Data.
	Registration UserRegistration

	// Public Name displayed at the Website.
	PublicName string
}

func ValidateUserPublicName(
	publicName string,
) (err error) {
	if (len(publicName) < UserPublicNameLengthMin) ||
		(len(publicName) > UserPublicNameLengthMax) {
		err = fmt.Errorf(ErrUserPublicName, publicName)
		return
	}
	return
}
