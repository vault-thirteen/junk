package app

import "fmt"

// Messages.
const (
	MsgAppStarted           = "Application has been started."
	MsgAppStopped           = "Application has been stopped."
	MsgSignalSigterm        = "'SIGTERM' Signal is received."
	MsgSignalSigint         = "'SIGINT' Signal is received."
	MsgHttpServerStop       = "HTTP Server has been stopped"
	MsgStorageConnection    = "Storage has been connected"
	MsgStorageCheck         = "Storage has been checked"
	MsgStorageDisconnection = "Storage has been disconnected"
	MsgIsEnabled            = "is enabled"
	MsgIsDisabled           = "is disabled"
)

// Message Formats.
const (
	MsgHttpServerStart = "HTTP Server Start (TLS %v)"
)

// Composes a Message for the Start of an HTTP Server.
func makeMsgHttpServerStart(
	isTlsEnabled bool,
) (msg string) {
	if isTlsEnabled {
		msg = fmt.Sprintf(MsgHttpServerStart, MsgIsEnabled)
	} else {
		msg = fmt.Sprintf(MsgHttpServerStart, MsgIsDisabled)
	}
	return
}
