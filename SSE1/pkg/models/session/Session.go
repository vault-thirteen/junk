package session

import (
	"database/sql"
	"time"
)

// Session.
type Session struct {

	// Basic Attributes...

	// Session Id.
	Id uint

	// Session User.
	User SessionUser

	// Time Points...

	// Time of the Session's Start.
	StartTime time.Time

	// Time of the last Activity on the Session.
	LastAccessTime time.Time

	// Time when the Session was closed (finished).
	EndTime sql.NullTime

	// A unique Marker of a Session and its Hash Sum.
	Marker     string
	MarkerHash string

	// Access Token Key which was created at the Moment of Token Generation.
	TokenKey string
}
