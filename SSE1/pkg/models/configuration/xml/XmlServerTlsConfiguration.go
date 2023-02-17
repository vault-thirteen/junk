package configuration

import "encoding/xml"

type XmlServerTlsConfiguration struct {
	XMLName xml.Name `xml:"TLS"`

	// Attributes.
	CertificateFile string `xml:"certificateFile,attr"`
	IsEnabled       bool   `xml:"isEnabled,attr"`
	KeyFile         string `xml:"keyFile,attr"`
}
