package simple

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"time"
)

type OptionalEventParameters struct {
	// ID of a user who initiated the event.
	UserId *cmb.Id `json:"userId"`

	// Time of the event.
	Time *time.Time `json:"time"`
}
