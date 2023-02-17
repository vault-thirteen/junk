package configuration

// Timeout Interval Settings for Spam-Requests.
type ServerStorageCoolDownPeriods struct {

	// Cool-down Interval for User Disabling Requests, Seconds.
	UserUnreg uint

	// Cool-down Interval for User Logging-in Requests, Seconds.
	UserLogIn uint
}
