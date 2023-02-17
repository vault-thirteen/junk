package configuration

import "encoding/xml"

type XmlServerHttpServerConfiguration struct {
	XMLName xml.Name `xml:"HttpServer"`

	// Attributes.
	Address            string `xml:"address,attr"`
	CookiePath         string `xml:"cookiePath,attr"`
	ShutdownTimeoutSec uint   `xml:"shutdownTimeoutSec,attr"`
	TokenHeader        string `xml:"tokenHeader,attr"`
}
