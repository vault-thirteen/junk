package server

import (
	"context"
	"fmt"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/avm"
	server2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	rc "github.com/vault-thirteen/RingCaptcha/server"
	rs "github.com/vault-thirteen/SimpleBB/pkg/RCS/settings"
)

type Server struct {
	settings *rs.Settings

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
	ssp         *avm.SSP

	// Captcha manager.
	captchaManager *rc.CaptchaManager

	// JSON-RPC server.
	js *jrm1.Processor
}

func NewServer(s base.ISettings) (srv *Server, err error) {
	stn := s.(*rs.Settings)

	err = stn.Check()
	if err != nil {
		return nil, err
	}

	srv = &Server{
		settings:      stn,
		listenDsn:     net.JoinHostPort(stn.HttpSettings.Host, strconv.FormatUint(uint64(stn.HttpSettings.Port), 10)),
		mustBeStopped: make(chan bool, server2.MustBeStoppedChannelSize),
		subRoutines:   new(sync.WaitGroup),
		mustStop:      new(atomic.Bool),
		httpErrors:    make(chan error, server2.HttpErrorsChannelSize),
		ssp:           avm.NewSSP(),
	}
	srv.mustStop.Store(false)

	// RPC server.
	err = srv.initRpc()
	if err != nil {
		return nil, err
	}

	err = srv.initCaptchaManager()
	if err != nil {
		return nil, err
	}

	// HTTP Server.
	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouter)),
	}

	return srv, nil
}

func (srv *Server) GetListenDsn() (dsn string) {
	return srv.listenDsn
}

func (srv *Server) GetCaptchaManagerListenDsn() (dsn string) {
	return srv.captchaManager.GetListenDsn()
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

	srv.subRoutines.Add(2)
	go srv.listenForHttpErrors()
	go srv.clearJunk()

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

	srv.subRoutines.Wait()

	err = srv.captchaManager.Stop()
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
		listenError = srv.httpServer.ListenAndServe()

		if (listenError != nil) && (listenError != http.ErrServerClosed) {
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

func (srv *Server) clearJunk() {
	defer srv.subRoutines.Done()

	var err error
	var imageTTL = int(srv.settings.CaptchaSettings.ImageTtlSec)

main_loop:
	for {
		// Check & sleep.
		for i := 0; i < imageTTL; i++ {
			time.Sleep(time.Second)

			if srv.mustStop.Load() {
				break main_loop
			}
		}

		// Work.
		err = srv.captchaManager.ClearJunk()
		if err != nil {
			log.Println(err)
			srv.mustBeStopped <- true
		}
	}

	log.Println(server2.MsgJunkCleanerHasStopped)
}

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	srv.js.ServeHTTP(rw, req)
}

func (srv *Server) initCaptchaManager() (err error) {
	srv.captchaManager, err = rc.NewCaptchaManager(
		srv.settings.CaptchaSettings.StoreImages,
		srv.settings.CaptchaSettings.ImagesFolder,
		srv.settings.CaptchaSettings.ImageWidth,
		srv.settings.CaptchaSettings.ImageHeight,
		srv.settings.CaptchaSettings.ImageTtlSec,
		srv.settings.CaptchaSettings.ClearImagesFolderAtStart,
		srv.settings.CaptchaSettings.UseHttpServerForImages,
		srv.settings.CaptchaSettings.HttpServerHost,
		srv.settings.CaptchaSettings.HttpServerPort,
		&srv.httpErrors,
		srv.settings.CaptchaSettings.HttpServerName,
		srv.settings.CaptchaSettings.IsCachingEnabled,
		srv.settings.CaptchaSettings.CacheSizeLimit,
		srv.settings.CaptchaSettings.CacheVolumeLimit,
		srv.settings.CaptchaSettings.CacheRecordTtl,
	)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) ReportStart() {
	fmt.Println(server2.MsgRpcHttpServer + srv.GetListenDsn())

	if srv.settings.CaptchaSettings.UseHttpServerForImages {
		fmt.Println(server2.MsgImagesHttpServer + srv.GetCaptchaManagerListenDsn())
	}
}

func (srv *Server) UseConstructor(stn base.ISettings) (base.IServer, error) {
	return NewServer(stn)
}
