package server

import (
	"context"
	"errors"
	"fmt"
	c "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/models/complex/IncidentManager"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/Scheduler"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/avm"
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
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/dbo"
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/km"
	as "github.com/vault-thirteen/SimpleBB/pkg/ACM/settings"
	rp "github.com/vault-thirteen/auxie/rpofs"
)

type Server struct {
	// Settings.
	settings *as.Settings

	// HTTPS server.
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

	// Verification code generator.
	vcg *rp.Generator

	// Clients for external services.
	gwmServiceClient  *cc.Client
	rcsServiceClient  *cc.Client
	smtpServiceClient *cc.Client

	// Generator of request IDs.
	ridg *rp.Generator

	// JWT key maker.
	jwtkm *km.KeyMaker

	// Incident manager.
	incidentManager derived2.IIncidentManager

	// Scheduler.
	scheduler *cm.Scheduler
}

func NewServer(s base.ISettings) (srv *Server, err error) {
	stn := s.(*as.Settings)

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
	sp := dbo.SystemParameters{
		PreSessionExpirationTime: srv.settings.SystemSettings.PreSessionExpirationTime,
	}

	srv.dbo = dbo.NewDatabaseObject(srv.settings.DbSettings, sp)

	err = srv.dbo.Init()
	if err != nil {
		return nil, err
	}

	// HTTPS Server.
	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouter)),
	}

	err = srv.initVerificationCodeGenerator()
	if err != nil {
		return nil, err
	}

	err = srv.initRequestIdGenerator()
	if err != nil {
		return nil, err
	}

	err = srv.initJwtKeyMaker()
	if err != nil {
		return nil, err
	}

	err = srv.createClientsForExternalServices()
	if err != nil {
		return nil, err
	}

	err = srv.initIncidentManager()
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

func (srv *Server) initVerificationCodeGenerator() (err error) {
	symbols := c.MakeSymbolsNumbersAndCapitalLatinLetters()

	srv.vcg, err = rp.NewGenerator(srv.settings.SystemSettings.VerificationCodeLength.AsInt(), symbols)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) createClientsForExternalServices() (err error) {
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

	// RCS module.
	{
		var rcsSCS = &cset.ServiceClientSettings{
			Schema:                      srv.settings.RcsSettings.Schema,
			Host:                        srv.settings.RcsSettings.Host,
			Port:                        srv.settings.RcsSettings.Port,
			Path:                        srv.settings.RcsSettings.Path,
			EnableSelfSignedCertificate: srv.settings.RcsSettings.EnableSelfSignedCertificate,
		}

		srv.rcsServiceClient, err = cc.NewClientWithSCS(rcsSCS, app.ServiceShortName_RCS)
		if err != nil {
			return err
		}
	}

	// SMTP module.
	{
		var smtpSCS = &cset.ServiceClientSettings{
			Schema:                      srv.settings.SmtpSettings.Schema,
			Host:                        srv.settings.SmtpSettings.Host,
			Port:                        srv.settings.SmtpSettings.Port,
			Path:                        srv.settings.SmtpSettings.Path,
			EnableSelfSignedCertificate: srv.settings.SmtpSettings.EnableSelfSignedCertificate,
		}

		srv.smtpServiceClient, err = cc.NewClientWithSCS(smtpSCS, app.ServiceShortName_SMTP)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *Server) pingClientsForExternalServices() (err error) {
	// GWM module [optional].
	{
		if srv.settings.SystemSettings.IsTableOfIncidentsUsed {
			err = srv.gwmServiceClient.Ping(true)
			if err != nil {
				return err
			}
		}
	}

	// RCS module.
	{
		err = srv.rcsServiceClient.Ping(true)
		if err != nil {
			return err
		}
	}

	// SMTP module.
	{
		err = srv.smtpServiceClient.Ping(true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *Server) initRequestIdGenerator() (err error) {
	symbols := c.MakeSymbolsNumbersAndCapitalLatinLetters()

	srv.ridg, err = rp.NewGenerator(srv.settings.SystemSettings.LogInRequestIdLength.AsInt(), symbols)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) initJwtKeyMaker() (err error) {
	srv.jwtkm, err = km.New(
		srv.settings.JWTSettings.SigningMethod,
		srv.settings.JWTSettings.PrivateKeyFilePath,
		srv.settings.JWTSettings.PublicKeyFilePath,
	)
	if err != nil {
		return err
	}

	return nil
}

// This method uses the GWM service client as an argument, thus it should be
// called after initialisation of all external service clients.
func (srv *Server) initIncidentManager() (err error) {
	srv.incidentManager = im.NewIncidentManager(srv.settings.SystemSettings.IsTableOfIncidentsUsed, srv.dbo, srv.gwmServiceClient, &srv.settings.SystemSettings.BlockTimePerIncident)

	return nil
}

func (srv *Server) initScheduler() {
	funcs60 := []simple.ScheduledFn{
		srv.clearPreRegUsersTable,
		srv.clearPasswordChangesTable,
		srv.clearEmailChangesTable,
		srv.clearSessions,
	}
	srv.scheduler = cm.NewScheduler(srv, funcs60, nil, nil)
}

func (srv *Server) ReportStart() {
	fmt.Println(server2.MsgHttpsServer + srv.GetListenDsn())
}

func (srv *Server) UseConstructor(stn base.ISettings) (base.IServer, error) {
	return NewServer(stn)
}
