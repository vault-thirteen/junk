package server

import (
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/dbo"
)

func (srv *Server) clearIPAddresses() (err error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	_, err = srv.dbo.GetPreparedStatementByIndex(dbo.DbPsid_ClearIPAddresses).Exec()
	if err != nil {
		return err
	}

	return nil
}
