package request

import (
	"errors"
	"fmt"
	"net"
)

// Settings.
const (
	MachineHostLengthMin             = 1
	MachineHostLengthMax             = 1024
	MachineBrowserUserAgentLengthMin = 3
	MachineBrowserUserAgentLengthMax = 4096
	NetworkIPType                    = "ip"
)

// Errors.
const (
	ErrfMachineIPAddress       = "Error in Machine.Host: '%v'."
	ErrMachineBrowserUserAgent = "Error in Machine.BrowserUserAgent."
)

type UserLogRequestMachine struct {
	Host             string
	BrowserUserAgent UserLogRequestMachineBrowserUserAgent
}

func ValidateMachineHost(
	host string,
) (err error) {
	if (len(host) < MachineHostLengthMin) ||
		(len(host) > MachineHostLengthMax) {
		err = fmt.Errorf(ErrfMachineIPAddress, host)
		return
	}
	_, err = net.ResolveIPAddr(NetworkIPType, host)
	if err != nil {
		return
	}
	return
}

func ValidateMachineBrowserUserAgent(
	bua string,
) (err error) {
	if (len(bua) < MachineBrowserUserAgentLengthMin) ||
		(len(bua) > MachineBrowserUserAgentLengthMax) {
		err = errors.New(ErrMachineBrowserUserAgent)
		return
	}
	return
}
