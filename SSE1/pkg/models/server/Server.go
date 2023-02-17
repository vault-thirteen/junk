package server

import (
	"time"

	"github.com/vault-thirteen/junk/SSE1/pkg/interfaces/logger"
)

// Server.
type Server struct {

	// Time of Start and Stop,
	// converted into UTC Time Zone.
	TimeOfStart time.Time
	TimeOfStop  time.Time

	// Logger.
	Logger logger.ILogger
}
