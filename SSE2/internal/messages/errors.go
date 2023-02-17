package messages

import "github.com/pkg/errors"

var (
	ErrIsAlreadyStarted = errors.New("is already started")
	ErrIsNotStarted     = errors.New("is not started")
)
