package rm

import (
	"net"
)

// Auth is authorisation data.
// This field must always be set for all requests except 'Ping'.
type Auth struct {
	// User's IP address, in a common textual format.
	// This field must always be set for all requests.
	UserIPA string `json:"userIPA"`

	// Parsed IP address of a user.	B = Byte array.
	// This field is used only in internal communications.
	UserIPAB net.IP `json:"-"`

	// Authorisation token.
	// This field is set for all requests of a logged user.
	Token string `json:"token"`
}
