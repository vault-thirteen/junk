package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vault-thirteen/junk/SSE1/pkg/helper/file"
	loggerHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/logger"
	"github.com/vault-thirteen/junk/SSE1/pkg/interfaces/logger"
	"github.com/vault-thirteen/junk/SSE1/pkg/interfaces/storage"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/buam"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/configuration"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/db/mysql"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/logger/builtin"
)

// Application Channels' Settings.
const (
	ErrorChannelSize = 256
	OssChannelSize   = 16
	QuitChannelSize  = 16
)

// Other hard-coded Settings.
const (
	TablesIniScriptsFolder = "Tables"
	SqlFileFullExt         = ".sql"
)

// Application.
type Application struct {

	// Settings.
	configuration *configuration.AppConfiguration

	// Logger.
	logger logger.ILogger

	// Channels.
	errorChannel chan error
	quitChannel  chan bool
	ossChannel   chan os.Signal

	// Network Servers and Data Sources.
	httpServer *http.Server
	httpRouter *httprouter.Router
	storage    storage.IStorage
	buam       *buam.BrowserUserAgentManager
}

// Application Constructor.
func NewApplication(
	configuration *configuration.AppConfiguration,
) (app *Application, err error) {
	app = &Application{}
	err = app.init(configuration)
	if err != nil {
		return
	}
	return
}

// Application Initialization.
func (app *Application) init(
	configuration *configuration.AppConfiguration,
) (err error) {
	app.configuration = configuration
	app.errorChannel = make(chan error, ErrorChannelSize)
	app.quitChannel = make(chan bool, QuitChannelSize)
	app.ossChannel = make(chan os.Signal, OssChannelSize)
	err = app.initLogger()
	if err != nil {
		return
	}
	go app.receiveErrors()
	go app.listenToOSSignals()
	err = app.initHttpServer()
	if err != nil {
		return
	}
	err = app.initStorage()
	if err != nil {
		return
	}
	err = app.initBUAM()
	if err != nil {
		return
	}
	return
}

// Initialization of a Logger.
func (app *Application) initLogger() (err error) {
	if app.configuration.Server.Logger.IsEnabled {
		switch app.configuration.Server.Logger.Type {
		case configuration.ServerLoggerTypeBuiltIn:
			app.logger = &builtin.BuiltInILogger{}
		default:
			err = errors.New(ErrUnsupportedLoggerType)
			return
		}
	}
	return
}

// Initialization of an HTTP Server.
func (app *Application) initHttpServer() (err error) {
	app.httpRouter = httprouter.New()
	err = app.initHttpRouter()
	if err != nil {
		return
	}
	app.httpServer = &http.Server{
		Addr:    app.configuration.Server.HttpServer.Address,
		Handler: app.httpRouter,
	}
	return
}

// Initialization of a Storage.
func (app *Application) initStorage() (err error) {
	err = app.initStorageScripts()
	if err != nil {
		return
	}
	switch app.configuration.Server.Storage.Type {
	case configuration.ServerStorageTypeMysql:
		app.storage, err = mysql.NewMysqlStorage(
			app.configuration.Server.Storage,
			app.logger,
		)
		if err != nil {
			return
		}
	default:
		err = errors.New(ErrUnsupportedStorageType)
		return
	}
	return
}

// Initialization of a Manager of Browser User Agent Records.
func (app *Application) initBUAM() (err error) {
	app.buam, err = buam.NewBrowserUserAgentManager(app.storage)
	if err != nil {
		return
	}
	return
}

// Runs the Storage Initialization Scripts.
// Each Script contains an SQL Code for a single Database Table.
func (app *Application) initStorageScripts() (err error) {
	var list = make(map[string]string)
	// Read Initialization Scripts for each Database Table.
	var fileContentsts string
	var filePath string
	for _, ts := range app.configuration.Server.Storage.TableSettings {
		filePath = filepath.Join(
			app.configuration.Server.Storage.InitializationScripts.Folder,
			TablesIniScriptsFolder,
			ts.TableName+SqlFileFullExt,
		)
		fileContentsts, err = file.GetTextFileContents(filePath)
		if err != nil {
			return
		}
		list[ts.TableName] = fileContentsts
	}
	app.configuration.Server.Storage.InitializationScripts.TableScripts = list
	return
}

// Runs the Application.
func (app *Application) Run() (err error) {
	err = app.start()
	if err != nil {
		return
	}
	app.log(MsgAppStarted)

	// Wait for a Quit-Request.
	var request bool
	for request = range app.quitChannel {
		if request == true {
			err = app.shutdown()
			if err != nil {
				app.errorChannel <- err
			} else {
				break
			}
		}
	}
	app.log(MsgAppStopped)
	return
}

// Starts the Application.
func (app *Application) start() (err error) {
	err = app.connectStorage()
	if err != nil {
		return
	}
	err = app.checkStorage()
	if err != nil {
		return
	}
	err = app.startHttpServer()
	if err != nil {
		return
	}
	return
}

// Starts the HTTP Server.
func (app *Application) startHttpServer() (err error) {
	go func() {
		var msg = makeMsgHttpServerStart(app.configuration.Server.TLS.IsEnabled)
		app.log(msg)
		var listenError error
		if app.configuration.Server.TLS.IsEnabled {
			listenError = app.httpServer.ListenAndServeTLS(
				app.configuration.Server.TLS.CertificateFile,
				app.configuration.Server.TLS.KeyFile,
			)
		} else {
			listenError = app.httpServer.ListenAndServe()
		}
		if listenError != nil {
			if listenError != http.ErrServerClosed {
				app.errorChannel <- listenError
			}
		}
	}()
	return
}

// Connects the Storage.
func (app *Application) connectStorage() (err error) {
	err = app.storage.Connect()
	if err != nil {
		return
	}
	app.log(MsgStorageConnection)
	return
}

// Checks Integrity of the Storage.
func (app *Application) checkStorage() (err error) {
	err = app.storage.Check()
	if err != nil {
		return
	}
	app.log(MsgStorageCheck)
	return
}

// Logs a Message.
func (app *Application) log(
	message string,
) {
	if app.configuration.Server.Logger.IsEnabled {
		loggerHelper.UseLogger(app.logger, message)
	}
}

// Runs the Manager which listens to Signals from an Operating System.
func (app *Application) listenToOSSignals() {
	signal.Notify(
		app.ossChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	var sig os.Signal
	for sig = range app.ossChannel {
		switch sig {
		case syscall.SIGTERM:
			app.log(MsgSignalSigterm)
			app.quitChannel <- true
		case syscall.SIGINT:
			app.log(MsgSignalSigint)
			app.quitChannel <- true
		}
	}
}

// Runs the Manager which listens to Application Errors.
func (app *Application) receiveErrors() {
	var err error
	for err = range app.errorChannel {
		app.log(err.Error())
	}
}

// Stops the Application gracefully.
func (app *Application) shutdown() (err error) {
	err = app.stopHttpServer()
	if err != nil {
		return
	}
	err = app.disconnectStorage()
	if err != nil {
		return
	}
	return
}

// Stops the HTTP Server gracefully.
func (app *Application) stopHttpServer() (err error) {
	var ctx context.Context
	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(app.configuration.Server.HttpServer.ShutdownTimeoutSec),
	)
	defer cancelFunc()
	err = app.httpServer.Shutdown(ctx)
	if err != nil {
		return
	}
	app.log(MsgHttpServerStop)
	return
}

// Disconnects the Storage.
func (app *Application) disconnectStorage() (err error) {
	err = app.storage.Disconnect()
	if err != nil {
		return
	}
	app.log(MsgStorageDisconnection)
	return
}
