package server

import (
	"errors"
	"net"
)

// Auxiliary functions used in RPC functions.

const (
	ErrFDatabaseNetwork = "database network error: %v" // Template.

)

// IsNetworkError checks if an error is a network error.
func IsNetworkError(err error) (isNetworkError bool) {
	var nerr net.Error
	return errors.As(err, &nerr)
}
