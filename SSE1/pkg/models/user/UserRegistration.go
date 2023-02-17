package user

import (
	"fmt"
	"time"
)

// Settings.
const (
	UserRegistrationSecretCodeLengthMin = 3
	UserRegistrationSecretCodeLengthMax = 64
)

// Errors.
const (
	ErrUserRegistrationSecretCode = "Error in UserRegistration.SecretCode: '%v'."
)

// User Registration Data.
type UserRegistration struct {

	// Time of User Registration.
	Time time.Time

	// Secret Code used for User Disabling.
	SecretCode string
}

func ValidateUserRegistrationSecretCode(
	secretCode string,
) (err error) {
	if (len(secretCode) < UserRegistrationSecretCodeLengthMin) ||
		(len(secretCode) > UserRegistrationSecretCodeLengthMax) {
		err = fmt.Errorf(ErrUserRegistrationSecretCode, secretCode)
		return
	}
	return
}
