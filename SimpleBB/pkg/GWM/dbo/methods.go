package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	dbo2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/dbo"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
	"net"
)

func (dbo *DatabaseObject) CountBlocksByIPAddress(ipa net.IP) (n cmb.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountBlocksByIPAddress).QueryRow(ipa)

	n, err = cms.NewNonNullValueFromScannableSource[cmb.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) InsertBlock(ipa net.IP, durationSec cmb.Count) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_AddBlock).Exec(ipa, durationSec)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) IncreaseBlockDuration(ipa net.IP, deltaDurationSec cmb.Count) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_IncreaseBlockDuration).Exec(deltaDurationSec, ipa)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}
