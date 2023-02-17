package configuration

import "encoding/xml"

type XmlServerLoggerConfiguration struct {
	XMLName xml.Name `xml:"Logger"`

	// Attributes.
	IsEnabled bool   `xml:"isEnabled,attr"`
	Type      string `xml:"type,attr"`
}
