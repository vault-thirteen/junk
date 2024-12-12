package dbo

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// SystemParameters stores various system parameters not directly related to
// the database.
type SystemParameters struct {
	PreSessionExpirationTime cmb.Count
}
