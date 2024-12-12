package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/vault-thirteen/SimpleBB/pkg/NM/models/complex/IncidentManager"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/DKey"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/Scheduler"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/avm"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	server2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	cset "github.com/vault-thirteen/SimpleBB/pkg/common/models/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/SimpleBB/pkg/NM/dbo"
	ns "github.com/vault-thirteen/SimpleBB/pkg/NM/settings"
)

type Server struct {
	// Settings.
	settings *ns.Settings

	// HTTP server.
	listenDsn  string
	httpServer *http.Server

	// Channel for an external controller. When a message comes from this
	// channel, a controller must stop this server. The server does not stop
	// itself.
	mustBeStopped chan bool

	// Internal control structures.
	subRoutines *sync.WaitGroup
	mustStop    *atomic.Bool
	httpErrors  chan error
	dbErrors    *chan error
	ssp         *avm.SSP

	// Database Object.
	dbo *dbo.DatabaseObject

	// JSON-RPC server.
	js *jrm1.Processor

	// Clients for external services.
	acmServiceClient *cc.Client
	gwmServiceClient *cc.Client
	smServiceClient  *cc.Client

	// Incident manager.
	incidentManager *im.IncidentManager

	// Internal DKeys.
	dKeyI *dk.DKey

	// External DKeys.
	dKeyForSM *cmb.Text

	// Scheduler.
	scheduler *cm.Scheduler
}

func NewServer(s base.ISettings) (srv *Server, err error) {
	stn := s.(*ns.Settings)

	err = stn.Check()
	if err != nil {
		return nil, err
	}

	dbErrorsChannel := make(chan error, server2.DbErrorsChannelSize)

	srv = &Server{
		settings:      stn,
		listenDsn:     net.JoinHostPort(stn.HttpsSettings.Host, strconv.FormatUint(uint64(stn.HttpsSettings.Port), 10)),
		mustBeStopped: make(chan bool, server2.MustBeStoppedChannelSize),
		subRoutines:   new(sync.WaitGroup),
		mustStop:      new(atomic.Bool),
		httpErrors:    make(chan error, server2.HttpErrorsChannelSize),
		dbErrors:      &dbErrorsChannel,
		ssp:           avm.NewSSP(),
	}
	srv.mustStop.Store(false)

	// RPC server.
	err = srv.initRpc()
	if err != nil {
		return nil, err
	}

	// Database.
	srv.dbo = dbo.NewDatabaseObject(srv.settings.DbSettings)

	err = srv.dbo.Init()
	if err != nil {
		return nil, err
	}

	// HTTP Server.
	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouter)),
	}

	err = srv.createClientsForExternalServices()
	if err != nil {
		return nil, err
	}

	err = srv.initIncidentManager()
	if err != nil {
		return nil, err
	}

	err = srv.initKeys()
	if err != nil {
		return nil, err
	}

	srv.initScheduler()

	return srv, nil
}

func (srv *Server) GetListenDsn() (dsn string) {
	return srv.listenDsn
}

func (srv *Server) GetStopChannel() *chan bool {
	return &srv.mustBeStopped
}

func (srv *Server) Start() (err error) {
	srv.ssp.Lock()
	defer srv.ssp.Unlock()

	err = srv.ssp.BeginStart()
	if err != nil {
		return err
	}

	srv.startHttpServer()

	srv.subRoutines.Add(3)
	go srv.listenForHttpErrors()
	go srv.listenForDbErrors()
	go srv.scheduler.Run()

	err = srv.incidentManager.Start()
	if err != nil {
		return err
	}

	err = srv.pingClientsForExternalServices()
	if err != nil {
		return err
	}

	err = srv.synchroniseModules(true)
	if err != nil {
		return err
	}

	srv.ssp.CompleteStart()

	return nil
}

func (srv *Server) Stop() (err error) {
	srv.ssp.Lock()
	defer srv.ssp.Unlock()

	err = srv.ssp.BeginStop()
	if err != nil {
		return err
	}

	srv.mustStop.Store(true)

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	err = srv.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	close(srv.httpErrors)
	close(*srv.dbErrors)

	srv.subRoutines.Wait()

	err = srv.incidentManager.Stop()
	if err != nil {
		return err
	}

	err = srv.dbo.Fin()
	if err != nil {
		return err
	}

	srv.ssp.CompleteStop()

	return nil
}

func (srv *Server) GetSubRoutinesWG() *sync.WaitGroup {
	return srv.subRoutines
}

func (srv *Server) GetMustStopAB() *atomic.Bool {
	return srv.mustStop
}

func (srv *Server) startHttpServer() {
	go func() {
		var listenError error
		listenError = srv.httpServer.ListenAndServeTLS(srv.settings.HttpsSettings.CertFile, srv.settings.HttpsSettings.KeyFile)

		if (listenError != nil) && (!errors.Is(listenError, http.ErrServerClosed)) {
			srv.httpErrors <- listenError
		}
	}()
}

func (srv *Server) listenForHttpErrors() {
	defer srv.subRoutines.Done()

	for err := range srv.httpErrors {
		log.Println(server2.MsgServerError + err.Error())
		srv.mustBeStopped <- true
	}

	log.Println(server2.MsgHttpErrorListenerHasStopped)
}

func (srv *Server) listenForDbErrors() {
	defer srv.subRoutines.Done()

	var err error
	for dbErr := range *srv.dbErrors {
		// When a network error occurs, it may be followed by numerous other
		// errors. If we try to fix each of them, we can make a flood. So,
		// we make a smart thing here.

		// 1. Ensure that the problem still exists.
		err = srv.dbo.ProbeDb()
		if err == nil {
			// Network is now fine. Ignore the error.
			continue
		}

		// 2. Log the error and try to reconnect.
		log.Println(server2.MsgDatabaseNetworkError + dbErr.Error())

		for {
			log.Println(server2.MsgReconnectingDatabase)
			// While we have prepared statements,
			// the simple reconnect will not help.
			err = srv.dbo.Init()
			if err != nil {
				// Network is still bad.
				log.Println(server2.MsgReconnectionHasFailed + err.Error())
			} else {
				log.Println(server2.MsgConnectionToDatabaseWasRestored)
				break
			}

			time.Sleep(time.Second * server2.DbReconnectCoolDownPeriodSec)
		}
	}

	log.Println(server2.MsgDbNetworkErrorListenerHasStopped)
}

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	srv.js.ServeHTTP(rw, req)
}

func (srv *Server) createClientsForExternalServices() (err error) {
	// ACM module.
	{
		var acmSCS = &cset.ServiceClientSettings{
			Schema:                      srv.settings.AcmSettings.Schema,
			Host:                        srv.settings.AcmSettings.Host,
			Port:                        srv.settings.AcmSettings.Port,
			Path:                        srv.settings.AcmSettings.Path,
			EnableSelfSignedCertificate: srv.settings.AcmSettings.EnableSelfSignedCertificate,
		}

		srv.acmServiceClient, err = cc.NewClientWithSCS(acmSCS, app.ServiceShortName_ACM)
		if err != nil {
			return err
		}
	}

	// GWM module [optional].
	{
		if !srv.settings.SystemSettings.IsTableOfIncidentsUsed {
			fmt.Println(server2.MsgIncidentsTableIsDisabled)
		} else {
			fmt.Println(server2.MsgIncidentsTableIsEnabled)

			var gwmSCS = &cset.ServiceClientSettings{
				Schema:                      srv.settings.GwmSettings.Schema,
				Host:                        srv.settings.GwmSettings.Host,
				Port:                        srv.settings.GwmSettings.Port,
				Path:                        srv.settings.GwmSettings.Path,
				EnableSelfSignedCertificate: srv.settings.GwmSettings.EnableSelfSignedCertificate,
			}

			srv.gwmServiceClient, err = cc.NewClientWithSCS(gwmSCS, app.ServiceShortName_GWM)
			if err != nil {
				return err
			}
		}
	}

	// SM module.
	{
		var smSCS = &cset.ServiceClientSettings{
			Schema:                      srv.settings.SmSettings.Schema,
			Host:                        srv.settings.SmSettings.Host,
			Port:                        srv.settings.SmSettings.Port,
			Path:                        srv.settings.SmSettings.Path,
			EnableSelfSignedCertificate: srv.settings.SmSettings.EnableSelfSignedCertificate,
		}

		srv.smServiceClient, err = cc.NewClientWithSCS(smSCS, app.ServiceShortName_SM)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *Server) pingClientsForExternalServices() (err error) {
	// ACM module.
	{
		err = srv.acmServiceClient.Ping(true)
		if err != nil {
			return err
		}
	}

	// GWM module.
	{
		err = srv.gwmServiceClient.Ping(true)
		if err != nil {
			return err
		}
	}

	// SM module.
	{
		err = srv.smServiceClient.Ping(true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *Server) synchroniseModules(verbose bool) (err error) {
	// SM module.
	{
		if verbose {
			fmt.Print(fmt.Sprintf(server2.MsgFSynchronisingWithModule, app.ServiceShortName_SM))
		}

		var re *jrm1.RpcError
		srv.dKeyForSM, re = srv.getDKeyForSM()
		if re != nil {
			return re.AsError()
		}

		if verbose {
			fmt.Println(server2.MsgOK)
		}
	}

	return nil
}

// This method uses the GWM service client as an argument, thus it should be
// called after initialisation of all external service clients.
func (srv *Server) initIncidentManager() (err error) {
	srv.incidentManager = im.NewIncidentManager(srv.settings.SystemSettings.IsTableOfIncidentsUsed.AsBool(), srv.dbo, srv.gwmServiceClient, &srv.settings.SystemSettings.BlockTimePerIncident)

	return nil
}

func (srv *Server) initScheduler() {
	funcs60 := []simple.ScheduledFn{
		srv.clearNotifications,
	}
	srv.scheduler = cm.NewScheduler(srv, funcs60, nil, nil)
}

func (srv *Server) initKeys() (err error) {
	srv.dKeyI, err = dk.NewDKey(int(srv.settings.SystemSettings.DKeySize))
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) ReportStart() {
	fmt.Println(server2.MsgHttpsServer + srv.GetListenDsn())
}

func (srv *Server) UseConstructor(stn base.ISettings) (base.IServer, error) {
	return NewServer(stn)
}
