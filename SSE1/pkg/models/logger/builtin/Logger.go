package builtin

import (
	"log"
)

// A built-in Logger.
// Implements the 'ILogger' Interface.
type BuiltInILogger struct {
}

// Logs the Message.
func (l *BuiltInILogger) Log(
	message string,
) (err error) {
	log.Println(message)
	return
}
