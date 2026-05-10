package avm

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	ErrAlreadyStarted = "already started"
	ErrAlreadyStopped = "already stopped"
)

// SSP is a Start-Stop Protector.
// It protects anything from multiple starts, multiple stops and simultaneous
// start-stops. This is a helper object which provides convenient access to
// functions used to protect starting and stopping procedures from stupid users.
type SSP struct {
	m         sync.Mutex
	isStarted atomic.Bool
}

func NewSSP() (ssp *SSP) {
	ssp = &SSP{}
	ssp.isStarted.Store(false)
	return ssp
}

func (ssp *SSP) Lock() {
	ssp.m.Lock()
}

func (ssp *SSP) Unlock() {
	ssp.m.Unlock()
}

func (ssp *SSP) BeginStart() (err error) {
	if ssp.isStarted.Load() {
		return errors.New(ErrAlreadyStarted)
	}
	return nil
}

func (ssp *SSP) CompleteStart() {
	ssp.isStarted.Store(true)
}

func (ssp *SSP) BeginStop() (err error) {
	if !ssp.isStarted.Load() {
		return errors.New(ErrAlreadyStopped)
	}
	return nil
}

func (ssp *SSP) CompleteStop() {
	ssp.isStarted.Store(false)
}
