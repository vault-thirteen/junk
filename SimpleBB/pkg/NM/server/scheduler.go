package server

import (
	"time"

	"github.com/vault-thirteen/SimpleBB/pkg/NM/dbo"
)

func (srv *Server) clearNotifications() (err error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	timeBorder := time.Now().Add(-time.Duration(srv.settings.SystemSettings.NotificationTtl) * time.Second)

	_, err = srv.dbo.GetPreparedStatementByIndex(dbo.DbPsid_ClearNotifications).Exec(timeBorder)
	if err != nil {
		return err
	}

	return nil
}
