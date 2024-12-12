package server

import (
	"fmt"
	server2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"log"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
)

// Auxiliary functions used in RPC functions.

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if err == nil {
		return
	}

	if srv.settings.GetSystemSettings().GetIsDebugMode() {
		log.Println(err)
	}
}

// processDatabaseError processes a database error.
func (srv *Server) processDatabaseError(err error) {
	if err == nil {
		return
	}

	if server2.IsNetworkError(err) {
		log.Println(fmt.Sprintf(server2.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return
}

// databaseError processes the database error and returns an RPC error.
func (srv *Server) databaseError(err error) (re *jrm1.RpcError) {
	srv.processDatabaseError(err)
	return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Database, server2.RpcErrorMsg_Database, err)
}
