package models

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type Subscription struct {
	// ID of a thread to which a user is subscribed.
	ThreadId cmb.Id `json:"threadId"`

	// ID of a subscribed user.
	UserId cmb.Id `json:"userId"`
}
