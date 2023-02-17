package configuration

import "encoding/xml"

type XmlServerAccessSessionConfiguration struct {
	XMLName xml.Name `xml:"Session"`

	// Children.
	IdleSessionTimeoutSec uint `xml:"idleSessionTimeoutSec,attr"`
}
