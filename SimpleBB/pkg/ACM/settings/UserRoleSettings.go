package s

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

// UserRoleSettings are settings for special user roles.
type UserRoleSettings struct {
	// List of IDs of users having a moderator role.
	ModeratorIds []cmb.Id `json:"moderatorIds"`

	// List of IDs of users having an administrator role.
	AdministratorIds []cmb.Id `json:"administratorIds"`
}

func (s UserRoleSettings) Check() (err error) {
	return nil
}
