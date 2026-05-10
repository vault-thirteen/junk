package net

import (
	"errors"
	"net"
	"strings"
)

const (
	ErrParseError             = "IP address parse error"
	ErrNotEnoughDataInAddress = "not enough data in address"
)

// ParseIPA parses a string into an IP address and returns an error on error.
// Golang's built-in parser does not return an error! What a shame.
func ParseIPA(s string) (ipa net.IP, err error) {
	ipa = net.ParseIP(s)

	if ipa == nil {
		return nil, errors.New(ErrParseError)
	}

	return ipa, nil
}

// SplitHostPort splits address into host and port and returns an error on
// error. While Go language does not have this very basic function, we are
// re-inventing the wheel again and again.
func SplitHostPort(addr string) (host, port string, err error) {
	parts := strings.Split(addr, ":")

	if len(parts) != 2 {
		return "", "", errors.New(ErrNotEnoughDataInAddress)
	}

	return parts[0], parts[1], nil
}

// SplitUrlPath splits an URL path into non-empty parts. This method is opposed
// to the standard 'split' method which allows empty parts.
func SplitUrlPath(path string) (parts []string) {
	p := strings.Split(path, "/")
	parts = make([]string, 0, len(p))
	for _, x := range p {
		if len(x) > 0 {
			parts = append(parts, x)
		}
	}

	return parts
}
