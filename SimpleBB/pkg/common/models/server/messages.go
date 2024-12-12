package server

// Ordinary messages.
const (
	MsgOK                               = "OK"
	MsgUsingDefaultConfigurationFile    = "Using the default configuration file"
	MsgServerIsStarting                 = "Server is starting ..."
	MsgHttpServer                       = "HTTP Server: "
	MsgHttpsServer                      = "HTTPS Server: "
	MsgRpcHttpServer                    = "RPC HTTP Server: "
	MsgImagesHttpServer                 = "Images HTTP Server: "
	MsgServerIsStopping                 = "Stopping the server ..."
	MsgServerIsStopped                  = "Server was stopped"
	MsgQuitSignalIsReceived             = "Quit signal from OS has been received: "
	MsgHttpErrorListenerHasStopped      = "HTTP error listener has stopped"
	MsgDbNetworkErrorListenerHasStopped = "DB network error listener has stopped"
	MsgEnterDatabasePassword            = "Enter the database password:"
	MsgEnterSmtpPassword                = "Enter the SMTP password:"
	MsgConnectingToDatabase             = "Connecting to database ..."
	MsgReconnectingDatabase             = "Reconnecting database ..."
	MsgReconnectionHasFailed            = "Reconnection has failed: "
	MsgConnectionToDatabaseWasRestored  = "Connection to database was restored"
	MsgSchedulerHasStopped              = "Scheduler has stopped"
	MsgJunkCleanerHasStopped            = "Junk cleaner has stopped"
	MsgIncidentManagerHasStopped        = "Incident manager has stopped"
	MsgIncidentsTableIsEnabled          = "Incidents table is enabled"
	MsgIncidentsTableIsDisabled         = "Incidents table is disabled"
	MsgFirewallIsEnabled                = "Firewall is enabled"
	MsgFirewallIsDisabled               = "Firewall is disabled"
	MsgPingAttempt                      = "."
	MsgDatabaseConsistencyCheck         = "Database consistency check ..."
)

// Error messages (simple).
const (
	MsgDatabaseNetworkError           = "Database network error: "
	MsgServerError                    = "Server error: "
	MsgSystemSettingError             = "Error in system setting"
	MsgSmtpSettingError               = "Error in SMTP module setting"
	MsgMessageSettingError            = "Error in message setting"
	MsgCaptchaServiceSettingError     = "Error in captcha service setting"
	MsgCaptchaImageServerSettingError = "Error in captcha image server setting"
	MsgJwtSettingError                = "Error in JWT setting"
)

// Templates for messages and errors.
const (
	MsgFTableIsNotFound            = "Table is not found: %s."
	MsgFInitialisingDatabaseTable  = "Initialising database table: %s."
	MsgFPingingModule              = "Pinging the %s module ..."
	MsgFModuleIsBroken             = "%s module is broken"
	MsgFServiceClientSettingsError = "%s service client settings error: %s"
	MsgFSynchronisingWithModule    = "Synchronising with %s module ..."
)
