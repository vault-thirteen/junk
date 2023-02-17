package logger

// Logger Interface.
type ILogger interface {

	// Logs the Message.
	Log(message string) error
}
