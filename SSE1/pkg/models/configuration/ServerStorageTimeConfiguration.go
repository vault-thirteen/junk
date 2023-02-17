package configuration

import "time"

type ServerStorageTimeConfiguration struct {
	Format string
	Zone   *time.Location
}
