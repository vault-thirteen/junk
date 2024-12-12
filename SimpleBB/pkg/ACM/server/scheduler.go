package server

import (
	"time"

	"github.com/vault-thirteen/SimpleBB/pkg/ACM/dbo"
)

func (srv *Server) clearPreRegUsersTable() (err error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.PreRegUserExpirationTime) * time.Second)

	_, err = srv.dbo.GetPreparedStatementByIndex(dbo.DbPsid_ClearPreRegUsersTable).Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) clearPasswordChangesTable() (err error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.PasswordChangeExpirationTime) * time.Second)

	_, err = srv.dbo.GetPreparedStatementByIndex(dbo.DbPsid_ClearPasswordChangesTable).Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) clearEmailChangesTable() (err error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.EmailChangeExpirationTime) * time.Second)

	_, err = srv.dbo.GetPreparedStatementByIndex(dbo.DbPsid_ClearEmailChangesTable).Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) clearSessions() (err error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.SessionMaxDuration) * time.Second)

	_, err = srv.dbo.GetPreparedStatementByIndex(dbo.DbPsid_ClearSessions).Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}
