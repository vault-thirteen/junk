package session

// Session User.
type SessionUser struct {

	// User Id.
	Id uint

	// Host Name of the User's Machine.
	Host string

	// 'User Agent' Field of the User's Internet Browser.
	BrowserUserAgent string

	// Internal Id of a 'User Agent' Field Value of the User's Internet Browser.
	BrowserUserAgentId uint
}
