package configuration

import "time"

type ServerTimezone struct {

	// Name of the Time Zone.
	Name string

	// Offset in Hours from Greenwich.
	OffsetHrs int8

	// Time Zone Location.
	location *time.Location
}
