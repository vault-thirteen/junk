package net

import "strings"

// Reads the Host from the Address String.
func GetAddressHost(
	address string,
) (host string) {
	const HostPortDelimiter = ":"
	parts := strings.Split(address, HostPortDelimiter)
	if len(parts) < 1 {
		return
	}
	return parts[0]
}
