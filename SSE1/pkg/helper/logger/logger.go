package logger

import (
	"fmt"
	"log"

	"github.com/vault-thirteen/junk/SSE1/pkg/interfaces/logger"
)

// Errors Formats.
const (
	ErrfLoggerError = "Logger Error: %v"
)

// A useful shorthand Function to call a Logger.
func UseLogger(
	logger logger.ILogger,
	message string,
) {
	var err = logger.Log(message)
	if err != nil {
		log.Println(message)
		log.Println(fmt.Errorf(ErrfLoggerError, err))
	}
}
