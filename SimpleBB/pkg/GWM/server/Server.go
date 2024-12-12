package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	a "github.com/vault-thirteen/SimpleBB/pkg/ACM/server"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/api"
	m "github.com/vault-thirteen/SimpleBB/pkg/MM/server"
	n "github.com/vault-thirteen/SimpleBB/pkg/NM/server"
	s "github.com/vault-thirteen/SimpleBB/pkg/SM/server"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/Scheduler"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/app"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/avm"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	ch "github.com/vault-thirteen/SimpleBB/pkg/common/models/http"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/models/net"
	server2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	cset "github.com/vault-thirteen/SimpleBB/pkg/common/models/settings"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/dbo"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	gs "github.com/vault-thirteen/SimpleBB/pkg/GWM/settings"
)

const ErrUrlIsTooShort = "URL is too short"

type Server struct {
	// Settings.
	settings gs.ISettings

	// HTTP server for internal requests.
	listenDsnInt  string
	httpServerInt *http.Server

	// HTTPS server for external requests.
	listenDsnExt     string
	httpServerExt    *http.Server
	apiFunctionNames []string
	apiHandlers      map[string]api.RequestHandler

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
	mmServiceClient  *cc.Client
	nmServiceClient  *cc.Client
	rcsProxy         *httputil.ReverseProxy
	smServiceClient  *cc.Client

	// Mapping of HTTP status codes by RPC error code for various services.
	commonHttpStatusCodesByRpcErrorCode map[int]int
	acmHttpStatusCodesByRpcErrorCode    map[int]int
	mmHttpStatusCodesByRpcErrorCode     map[int]int
	nmHttpStatusCodesByRpcErrorCode     map[int]int
	smHttpStatusCodesByRpcErrorCode     map[int]int

	// Public settings.
	publicSettingsFileData []byte // Cached file contents in JSON format.

	// Front End Data.
	frontEnd *models.FrontEndData

	// Scheduler.
	scheduler *cm.Scheduler
}

func NewServer(s base.ISettings) (srv *Server, err error) {
	stn := s.(gs.ISettings)

	err = stn.Check()
	if err != nil {
		return nil, err
	}

	dbErrorsChannel := make(chan error, server2.DbErrorsChannelSize)

	srv = &Server{
		settings:      stn,
		listenDsnInt:  net.JoinHostPort(stn.GetIntHttpSettings().Host, strconv.FormatUint(uint64(stn.GetIntHttpSettings().Port), 10)),
		listenDsnExt:  net.JoinHostPort(stn.GetExtHttpsSettings().Host, strconv.FormatUint(uint64(stn.GetExtHttpsSettings().Port), 10)),
		mustBeStopped: make(chan bool, server2.MustBeStoppedChannelSize),
		subRoutines:   new(sync.WaitGroup),
		mustStop:      new(atomic.Bool),
		httpErrors:    make(chan error, server2.HttpErrorsChannelSize),
		dbErrors:      &dbErrorsChannel,
		ssp:           avm.NewSSP(),
	}
	srv.mustStop.Store(false)

	if srv.settings.GetSystemSettings().GetIsFirewallUsed() {
		fmt.Println(server2.MsgFirewallIsEnabled)
	} else {
		fmt.Println(server2.MsgFirewallIsDisabled)
	}

	// RPC server.
	err = srv.initRpc()
	if err != nil {
		return nil, err
	}

	// Database.
	srv.dbo = dbo.NewDatabaseObject(srv.settings.GetDbSettings())

	err = srv.dbo.Init()
	if err != nil {
		return nil, err
	}

	// HTTP server for internal requests.
	srv.httpServerInt = &http.Server{
		Addr:    srv.listenDsnInt,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouterInt)),
	}

	// HTTPS server for external requests.
	srv.httpServerExt = &http.Server{
		Addr:    srv.listenDsnExt,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouterExt)),
	}

	err = srv.initApiFunctions()
	if err != nil {
		return nil, err
	}

	err = srv.createClientsForExternalServices()
	if err != nil {
		return nil, err
	}

	err = srv.initStatusCodeMapper()
	if err != nil {
		return nil, err
	}

	err = srv.initPublicSettings()
	if err != nil {
		return nil, err
	}

	if srv.settings.GetSystemSettings().GetIsFrontEndEnabled() {
		err = srv.initFrontEndData()
		if err != nil {
			return nil, err
		}
	}

	srv.initScheduler()

	return srv, nil
}

func (srv *Server) GetListenDsnInt() (dsn string) {
	return srv.listenDsnInt
}

func (srv *Server) GetListenDsnExt() (dsn string) {
	return srv.listenDsnExt
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

	srv.startHttpServerInt()
	srv.startHttpServerExt()

	srv.subRoutines.Add(3)
	go srv.listenForHttpErrors()
	go srv.listenForDbErrors()
	go srv.scheduler.Run()

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

	ctxInt, cfInt := context.WithTimeout(context.Background(), time.Minute)
	defer cfInt()
	err = srv.httpServerInt.Shutdown(ctxInt)
	if err != nil {
		return err
	}

	ctxExt, cfExt := context.WithTimeout(context.Background(), time.Minute)
	defer cfExt()
	err = srv.httpServerExt.Shutdown(ctxExt)
	if err != nil {
		return err
	}

	close(srv.httpErrors)
	close(*srv.dbErrors)

	srv.subRoutines.Wait()

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

func (srv *Server) startHttpServerInt() {
	go func() {
		var listenError error
		listenError = srv.httpServerInt.ListenAndServe()

		if (listenError != nil) && (!errors.Is(listenError, http.ErrServerClosed)) {
			srv.httpErrors <- listenError
		}
	}()
}

func (srv *Server) startHttpServerExt() {
	go func() {
		var listenError error
		listenError = srv.httpServerExt.ListenAndServeTLS(srv.settings.GetExtHttpsSettings().CertFile, srv.settings.GetExtHttpsSettings().KeyFile)

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

// HTTP router for internal requests.
func (srv *Server) httpRouterInt(rw http.ResponseWriter, req *http.Request) {
	srv.js.ServeHTTP(rw, req)
}

// HTTP router for external requests.
func (srv *Server) httpRouterExt(rw http.ResponseWriter, req *http.Request) {
	// Firewall (optional).
	var clientIPA simple.IPAS
	if srv.settings.GetSystemSettings().GetIsFirewallUsed() {
		var ok bool
		var err error
		ok, clientIPA, err = srv.isIPAddressAllowed(req)
		if err != nil {
			srv.processInternalServerError(rw, err)
			return
		}

		if !ok {
			srv.respondForbidden(rw)
			return
		}
	}

	isFrontEndEnabled := srv.settings.GetSystemSettings().GetIsFrontEndEnabled()

	if isFrontEndEnabled {
		switch simple.Path(req.URL.Path) {
		case gs.FrontEndRoot: // <- /
			srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.IndexHtmlPage)
			return

		case srv.frontEnd.FavIcon.UrlPath: // <- /favicon.png
			srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.FavIcon)
			return

		case gs.FrontEndAdminPath: // <- /admin
			srv.handleAdminFrontEnd(rw, req, clientIPA)
			return
		}
	}

	var err error
	urlParts := cn.SplitUrlPath(req.URL.Path)
	if len(urlParts) < 1 {
		err = errors.New(ErrUrlIsTooShort)
		srv.logError(err)
		return
	}
	category := urlParts[0]

	switch simple.Path(category) {
	case srv.settings.GetSystemSettings().GetApiFolder(): // <- /api
		srv.handlePublicApi(rw, req, clientIPA)
		return

	case srv.settings.GetSystemSettings().GetPublicSettingsFileName(): // <- /settings.json
		srv.handlePublicSettings(rw, req)
		return

	case srv.settings.GetSystemSettings().GetCaptchaFolder(): // <- /captcha
		srv.handleCaptcha(rw, req)
		return

	case srv.settings.GetSystemSettings().GetFrontEndStaticFilesFolder(): // <- /fe
		if !isFrontEndEnabled {
			srv.respondNotFound(rw)
			return
		}
		srv.handleFrontEnd(rw, req, clientIPA)
		return
	}

	srv.respondNotFound(rw)
	return
}

func (srv *Server) initStatusCodeMapper() (err error) {
	srv.commonHttpStatusCodesByRpcErrorCode = server2.GetMapOfHttpStatusCodesByRpcErrorCodes()
	srv.acmHttpStatusCodesByRpcErrorCode = a.GetMapOfHttpStatusCodesByRpcErrorCodes()
	srv.mmHttpStatusCodesByRpcErrorCode = m.GetMapOfHttpStatusCodesByRpcErrorCodes()
	srv.nmHttpStatusCodesByRpcErrorCode = n.GetMapOfHttpStatusCodesByRpcErrorCodes()
	srv.smHttpStatusCodesByRpcErrorCode = s.GetMapOfHttpStatusCodesByRpcErrorCodes()

	return nil
}

func (srv *Server) initPublicSettings() (err error) {
	// File with public settings is a special virtual file which lies in the
	// root folder. It contains useful settings for client applications.
	var publicSettings = &models.Settings{
		Version:                   srv.settings.GetSystemSettings().GetSettingsVersion(),
		ProductVersion:            cmb.Text(srv.settings.GetVersionInfo().ProgramVersionString()),
		SiteName:                  srv.settings.GetSystemSettings().GetSiteName(),
		SiteDomain:                srv.settings.GetSystemSettings().GetSiteDomain(),
		CaptchaFolder:             srv.settings.GetSystemSettings().GetCaptchaFolder(),
		SessionMaxDuration:        srv.settings.GetSystemSettings().GetSessionMaxDuration(),
		MessageEditTime:           srv.settings.GetSystemSettings().GetMessageEditTime(),
		PageSize:                  srv.settings.GetSystemSettings().GetPageSize(),
		ApiFolder:                 srv.settings.GetSystemSettings().GetApiFolder(),
		PublicSettingsFileName:    srv.settings.GetSystemSettings().GetPublicSettingsFileName(),
		IsFrontEndEnabled:         srv.settings.GetSystemSettings().GetIsFrontEndEnabled(),
		FrontEndStaticFilesFolder: srv.settings.GetSystemSettings().GetFrontEndStaticFilesFolder(),
		NotificationCountLimit:    srv.settings.GetSystemSettings().GetNotificationCountLimit(),
	}

	srv.publicSettingsFileData, err = json.Marshal(publicSettings)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) initFrontEndData() (err error) {
	frontendAssetsFolder := simple.NormalisePath(srv.settings.GetSystemSettings().GetFrontEndAssetsFolder())
	fep := gs.FrontEndRoot + srv.settings.GetSystemSettings().GetFrontEndStaticFilesFolder() + "/" // <- /fe/

	srv.frontEnd = &models.FrontEndData{}

	srv.frontEnd.AdminHtmlPage, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_AdminHtmlPage, ch.ContentType_HtmlPage, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.AdminJs, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_AdminJs, ch.ContentType_JavaScript, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.ApiJs, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_ApiJs, ch.ContentType_JavaScript, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.ArgonJs, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_ArgonJs, ch.ContentType_JavaScript, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.ArgonWasm, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_ArgonWasm, ch.ContentType_Wasm, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.BppJs, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_BppJs, ch.ContentType_JavaScript, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.CssStyles, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_CssStyles, ch.ContentType_CssStyle, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.FavIcon, err = models.NewFrontEndFileData(gs.FrontEndRoot, gs.FrontEndStaticFileName_FavIcon, ch.ContentType_PNG, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.IndexHtmlPage, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_IndexHtmlPage, ch.ContentType_HtmlPage, frontendAssetsFolder)
	if err != nil {
		return err
	}

	srv.frontEnd.LoaderScript, err = models.NewFrontEndFileData(fep, gs.FrontEndStaticFileName_LoaderScript, ch.ContentType_JavaScript, frontendAssetsFolder)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) createClientsForExternalServices() (err error) {
	// ACM module.
	{
		var acmSCS = &cset.ServiceClientSettings{
			Schema:                      srv.settings.GetAcmSettings().Schema,
			Host:                        srv.settings.GetAcmSettings().Host,
			Port:                        srv.settings.GetAcmSettings().Port,
			Path:                        srv.settings.GetAcmSettings().Path,
			EnableSelfSignedCertificate: srv.settings.GetAcmSettings().EnableSelfSignedCertificate,
		}

		srv.acmServiceClient, err = cc.NewClientWithSCS(acmSCS, app.ServiceShortName_ACM)
		if err != nil {
			return err
		}
	}

	// MM module.
	{
		var mmSCS = &cset.ServiceClientSettings{
			Schema:                      srv.settings.GetMmSettings().Schema,
			Host:                        srv.settings.GetMmSettings().Host,
			Port:                        srv.settings.GetMmSettings().Port,
			Path:                        srv.settings.GetMmSettings().Path,
			EnableSelfSignedCertificate: srv.settings.GetMmSettings().EnableSelfSignedCertificate,
		}

		srv.mmServiceClient, err = cc.NewClientWithSCS(mmSCS, app.ServiceShortName_MM)
		if err != nil {
			return err
		}
	}

	// NM module.
	{
		var nmSCS = &cset.ServiceClientSettings{
			Schema:                      srv.settings.GetNmSettings().Schema,
			Host:                        srv.settings.GetNmSettings().Host,
			Port:                        srv.settings.GetNmSettings().Port,
			Path:                        srv.settings.GetNmSettings().Path,
			EnableSelfSignedCertificate: srv.settings.GetNmSettings().EnableSelfSignedCertificate,
		}

		srv.nmServiceClient, err = cc.NewClientWithSCS(nmSCS, app.ServiceShortName_NM)
		if err != nil {
			return err
		}
	}

	// Proxy for captcha images (RCS).
	{
		targetAddr := cc.UrlSchemeHttp + "://" +
			net.JoinHostPort(
				srv.settings.GetSystemSettings().GetCaptchaImgServerHost(),
				strconv.FormatUint(uint64(srv.settings.GetSystemSettings().GetCaptchaImgServerPort()), 10),
			)
		var targetUrl *url.URL
		targetUrl, err = url.Parse(targetAddr)
		if err != nil {
			return err
		}

		srv.rcsProxy = httputil.NewSingleHostReverseProxy(targetUrl)
	}

	// SM module.
	{
		var smSCS = &cset.ServiceClientSettings{
			Schema:                      srv.settings.GetSmSettings().Schema,
			Host:                        srv.settings.GetSmSettings().Host,
			Port:                        srv.settings.GetSmSettings().Port,
			Path:                        srv.settings.GetSmSettings().Path,
			EnableSelfSignedCertificate: srv.settings.GetSmSettings().EnableSelfSignedCertificate,
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

	// MM module.
	{
		err = srv.mmServiceClient.Ping(true)
		if err != nil {
			return err
		}
	}

	// NM module.
	{
		err = srv.nmServiceClient.Ping(true)
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

func (srv *Server) initApiFunctions() (err error) {
	srv.apiFunctionNames = []string{
		// ACM.
		ApiFunctionName_RegisterUser,
		ApiFunctionName_GetListOfRegistrationsReadyForApproval,
		ApiFunctionName_RejectRegistrationRequest,
		ApiFunctionName_ApproveAndRegisterUser,
		ApiFunctionName_LogUserIn,
		ApiFunctionName_LogUserOut,
		ApiFunctionName_LogUserOutA,
		ApiFunctionName_GetListOfLoggedUsers,
		ApiFunctionName_GetListOfLoggedUsersOnPage,
		ApiFunctionName_GetListOfAllUsers,
		ApiFunctionName_GetListOfAllUsersOnPage,
		ApiFunctionName_IsUserLoggedIn,
		ApiFunctionName_ChangePassword,
		ApiFunctionName_ChangeEmail,
		ApiFunctionName_GetUserSession,
		ApiFunctionName_GetUserRoles,
		ApiFunctionName_ViewUserParameters,
		ApiFunctionName_SetUserRoleAuthor,
		ApiFunctionName_SetUserRoleWriter,
		ApiFunctionName_SetUserRoleReader,
		ApiFunctionName_GetSelfRoles,
		ApiFunctionName_BanUser,
		ApiFunctionName_UnbanUser,

		// MM.
		ApiFunctionName_AddSection,
		ApiFunctionName_ChangeSectionName,
		ApiFunctionName_ChangeSectionParent,
		ApiFunctionName_GetSection,
		ApiFunctionName_MoveSectionUp,
		ApiFunctionName_MoveSectionDown,
		ApiFunctionName_DeleteSection,
		ApiFunctionName_AddForum,
		ApiFunctionName_ChangeForumName,
		ApiFunctionName_ChangeForumSection,
		ApiFunctionName_GetForum,
		ApiFunctionName_MoveForumUp,
		ApiFunctionName_MoveForumDown,
		ApiFunctionName_DeleteForum,
		ApiFunctionName_AddThread,
		ApiFunctionName_ChangeThreadName,
		ApiFunctionName_ChangeThreadForum,
		ApiFunctionName_GetThread,
		ApiFunctionName_GetThreadNamesByIds,
		ApiFunctionName_MoveThreadUp,
		ApiFunctionName_MoveThreadDown,
		ApiFunctionName_DeleteThread,
		ApiFunctionName_AddMessage,
		ApiFunctionName_ChangeMessageText,
		ApiFunctionName_ChangeMessageThread,
		ApiFunctionName_GetMessage,
		ApiFunctionName_GetLatestMessageOfThread,
		ApiFunctionName_DeleteMessage,
		ApiFunctionName_ListThreadAndMessages,
		ApiFunctionName_ListThreadAndMessagesOnPage,
		ApiFunctionName_ListForumAndThreads,
		ApiFunctionName_ListForumAndThreadsOnPage,
		ApiFunctionName_ListSectionsAndForums,

		// NM.
		ApiFunctionName_AddNotification,
		ApiFunctionName_GetNotification,
		ApiFunctionName_GetNotifications,
		ApiFunctionName_GetNotificationsOnPage,
		ApiFunctionName_GetUnreadNotifications,
		ApiFunctionName_CountUnreadNotifications,
		ApiFunctionName_MarkNotificationAsRead,
		ApiFunctionName_DeleteNotification,
		ApiFunctionName_AddResource,
		ApiFunctionName_GetResource,
		ApiFunctionName_GetResourceValue,
		ApiFunctionName_GetListOfAllResourcesOnPage,
		ApiFunctionName_DeleteResource,

		// SM.
		ApiFunctionName_AddSubscription,
		ApiFunctionName_IsSelfSubscribed,
		ApiFunctionName_IsUserSubscribed,
		ApiFunctionName_CountSelfSubscriptions,
		ApiFunctionName_GetSelfSubscriptions,
		ApiFunctionName_GetSelfSubscriptionsOnPage,
		ApiFunctionName_GetUserSubscriptions,
		ApiFunctionName_GetUserSubscriptionsOnPage,
		ApiFunctionName_DeleteSelfSubscription,
		ApiFunctionName_DeleteSubscription,
	}

	srv.apiHandlers = map[string]api.RequestHandler{
		// ACM.
		ApiFunctionName_RegisterUser:                           srv.RegisterUser,
		ApiFunctionName_GetListOfRegistrationsReadyForApproval: srv.GetListOfRegistrationsReadyForApproval,
		ApiFunctionName_RejectRegistrationRequest:              srv.RejectRegistrationRequest,
		ApiFunctionName_ApproveAndRegisterUser:                 srv.ApproveAndRegisterUser,
		ApiFunctionName_LogUserIn:                              srv.LogUserIn,
		ApiFunctionName_LogUserOut:                             srv.LogUserOut,
		ApiFunctionName_LogUserOutA:                            srv.LogUserOutA,
		ApiFunctionName_GetListOfLoggedUsers:                   srv.GetListOfLoggedUsers,
		ApiFunctionName_GetListOfLoggedUsersOnPage:             srv.GetListOfLoggedUsersOnPage,
		ApiFunctionName_GetListOfAllUsers:                      srv.GetListOfAllUsers,
		ApiFunctionName_GetListOfAllUsersOnPage:                srv.GetListOfAllUsersOnPage,
		ApiFunctionName_IsUserLoggedIn:                         srv.IsUserLoggedIn,
		ApiFunctionName_ChangePassword:                         srv.ChangePassword,
		ApiFunctionName_ChangeEmail:                            srv.ChangeEmail,
		ApiFunctionName_GetUserSession:                         srv.GetUserSession,
		ApiFunctionName_GetUserName:                            srv.GetUserName,
		ApiFunctionName_GetUserRoles:                           srv.GetUserRoles,
		ApiFunctionName_ViewUserParameters:                     srv.ViewUserParameters,
		ApiFunctionName_SetUserRoleAuthor:                      srv.SetUserRoleAuthor,
		ApiFunctionName_SetUserRoleWriter:                      srv.SetUserRoleWriter,
		ApiFunctionName_SetUserRoleReader:                      srv.SetUserRoleReader,
		ApiFunctionName_GetSelfRoles:                           srv.GetSelfRoles,
		ApiFunctionName_BanUser:                                srv.BanUser,
		ApiFunctionName_UnbanUser:                              srv.UnbanUser,

		// MM.
		ApiFunctionName_AddSection:                  srv.AddSection,
		ApiFunctionName_ChangeSectionName:           srv.ChangeSectionName,
		ApiFunctionName_ChangeSectionParent:         srv.ChangeSectionParent,
		ApiFunctionName_GetSection:                  srv.GetSection,
		ApiFunctionName_MoveSectionUp:               srv.MoveSectionUp,
		ApiFunctionName_MoveSectionDown:             srv.MoveSectionDown,
		ApiFunctionName_DeleteSection:               srv.DeleteSection,
		ApiFunctionName_AddForum:                    srv.AddForum,
		ApiFunctionName_ChangeForumName:             srv.ChangeForumName,
		ApiFunctionName_ChangeForumSection:          srv.ChangeForumSection,
		ApiFunctionName_GetForum:                    srv.GetForum,
		ApiFunctionName_MoveForumUp:                 srv.MoveForumUp,
		ApiFunctionName_MoveForumDown:               srv.MoveForumDown,
		ApiFunctionName_DeleteForum:                 srv.DeleteForum,
		ApiFunctionName_AddThread:                   srv.AddThread,
		ApiFunctionName_ChangeThreadName:            srv.ChangeThreadName,
		ApiFunctionName_ChangeThreadForum:           srv.ChangeThreadForum,
		ApiFunctionName_GetThread:                   srv.GetThread,
		ApiFunctionName_GetThreadNamesByIds:         srv.GetThreadNamesByIds,
		ApiFunctionName_MoveThreadUp:                srv.MoveThreadUp,
		ApiFunctionName_MoveThreadDown:              srv.MoveThreadDown,
		ApiFunctionName_DeleteThread:                srv.DeleteThread,
		ApiFunctionName_AddMessage:                  srv.AddMessage,
		ApiFunctionName_ChangeMessageText:           srv.ChangeMessageText,
		ApiFunctionName_ChangeMessageThread:         srv.ChangeMessageThread,
		ApiFunctionName_GetMessage:                  srv.GetMessage,
		ApiFunctionName_GetLatestMessageOfThread:    srv.GetLatestMessageOfThread,
		ApiFunctionName_DeleteMessage:               srv.DeleteMessage,
		ApiFunctionName_ListThreadAndMessages:       srv.ListThreadAndMessages,
		ApiFunctionName_ListThreadAndMessagesOnPage: srv.ListThreadAndMessagesOnPage,
		ApiFunctionName_ListForumAndThreads:         srv.ListForumAndThreads,
		ApiFunctionName_ListForumAndThreadsOnPage:   srv.ListForumAndThreadsOnPage,
		ApiFunctionName_ListSectionsAndForums:       srv.ListSectionsAndForums,

		// NM.
		ApiFunctionName_AddNotification:             srv.AddNotification,
		ApiFunctionName_GetNotification:             srv.GetNotification,
		ApiFunctionName_GetNotifications:            srv.GetNotifications,
		ApiFunctionName_GetNotificationsOnPage:      srv.GetNotificationsOnPage,
		ApiFunctionName_GetUnreadNotifications:      srv.GetUnreadNotifications,
		ApiFunctionName_CountUnreadNotifications:    srv.CountUnreadNotifications,
		ApiFunctionName_MarkNotificationAsRead:      srv.MarkNotificationAsRead,
		ApiFunctionName_DeleteNotification:          srv.DeleteNotification,
		ApiFunctionName_AddResource:                 srv.AddResource,
		ApiFunctionName_GetResource:                 srv.GetResource,
		ApiFunctionName_GetResourceValue:            srv.GetResourceValue,
		ApiFunctionName_GetListOfAllResourcesOnPage: srv.GetListOfAllResourcesOnPage,
		ApiFunctionName_DeleteResource:              srv.DeleteResource,

		// SM.
		ApiFunctionName_AddSubscription:            srv.AddSubscription,
		ApiFunctionName_IsSelfSubscribed:           srv.IsSelfSubscribed,
		ApiFunctionName_IsUserSubscribed:           srv.IsUserSubscribed,
		ApiFunctionName_CountSelfSubscriptions:     srv.CountSelfSubscriptions,
		ApiFunctionName_GetSelfSubscriptions:       srv.GetSelfSubscriptions,
		ApiFunctionName_GetSelfSubscriptionsOnPage: srv.GetSelfSubscriptionsOnPage,
		ApiFunctionName_GetUserSubscriptions:       srv.GetUserSubscriptions,
		ApiFunctionName_GetUserSubscriptionsOnPage: srv.GetUserSubscriptionsOnPage,
		ApiFunctionName_DeleteSubscription:         srv.DeleteSubscription,
		ApiFunctionName_DeleteSelfSubscription:     srv.DeleteSelfSubscription,
	}

	return nil
}

func (srv *Server) initScheduler() {
	funcs60 := []simple.ScheduledFn{
		srv.clearIPAddresses,
	}
	srv.scheduler = cm.NewScheduler(srv, funcs60, nil, nil)
}

func (srv *Server) ReportStart() {
	fmt.Println(server2.MsgHttpServer + srv.GetListenDsnInt())
	fmt.Println(server2.MsgHttpsServer + srv.GetListenDsnExt())
}

func (srv *Server) UseConstructor(stn base.ISettings) (base.IServer, error) {
	return NewServer(stn)
}
