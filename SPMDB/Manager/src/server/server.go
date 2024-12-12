package server

import (
	"context"
	"log"
	"net/http"
	"spdb_manager/cli/printer"
	"spdb_manager/database"
	"spdb_manager/settings"
	"strconv"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	lock             sync.Mutex
	quitChan         *chan bool
	host             string
	port             uint16
	isEditorModeUsed bool
	httpServer       *http.Server
	httpRouter       *httprouter.Router
	dbFile           string
	dbo              *database.DbSimple

	// Folder in file system which contains GUI (front-end) scripts.
	guiFolder string
}

// New is server's constructor.
func New(
	quitChan *chan bool,
	host string,
	port uint16,
	isEditorModeUsed bool,
	dbFile string,
	guiFolder string,
) (srv *Server, err error) {
	srv = &Server{
		quitChan:         quitChan,
		host:             host,
		port:             port,
		isEditorModeUsed: isEditorModeUsed,
		dbFile:           dbFile,
		guiFolder:        guiFolder,
	}

	err = srv.init()
	if err != nil {
		return nil, err
	}

	return srv, nil
}

// init initializes the server.
func (srv *Server) init() (err error) {
	err = srv.initRouter()
	if err != nil {
		return err
	}

	srv.httpServer = &http.Server{
		Addr:    srv.host + ":" + strconv.FormatUint(uint64(srv.port), 10),
		Handler: srv.httpRouter,
	}

	srv.dbo, err = database.OpenSimple(srv.dbFile)
	if err != nil {
		return err
	}

	return nil
}

// Start starts the server.
func (srv *Server) Start() {
	go func() {
		derr := srv.httpServer.ListenAndServe()
		if derr != nil {
			srv.logError(derr)
		}
	}()
}

// logError logs the error if it occurs.
func (srv *Server) logError(err error) {
	if err != nil {
		printer.PrintError(err.Error())
	}
}

// stop gracefully stops the server.
func (srv *Server) stop() {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	defer func() {
		log.Println("Server Shutdown")
		*(srv.quitChan) <- true
	}()

	defer func() {
		derr := srv.dbo.Close()
		srv.logError(derr)
		log.Println("Database Shutdown")
	}()

	defer func() {
		ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*settings.ServerShutdownTimeoutSec)
		defer cancelFn()
		derr := srv.httpServer.Shutdown(ctx)
		srv.logError(derr)
		log.Println("HTTP Server Shutdown")
	}()
}
