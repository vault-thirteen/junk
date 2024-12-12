package ur

import (
	"database/sql"
	"errors"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type UserRoles struct {
	IsAdministrator cmb.Flag `json:"isAdministrator"`
	IsModerator     cmb.Flag `json:"isModerator"`
	IsAuthor        cmb.Flag `json:"isAuthor"`
	IsWriter        cmb.Flag `json:"isWriter"`
	IsReader        cmb.Flag `json:"isReader"`
	CanLogIn        cmb.Flag `json:"canLogIn"`
}

func NewUserRolesFromScannableSource(src cmi.IScannable) (ur *UserRoles, err error) {
	ur = &UserRoles{}

	err = src.Scan(
		&ur.IsAuthor,
		&ur.IsWriter,
		&ur.IsReader,
		&ur.CanLogIn,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ur, nil
}
