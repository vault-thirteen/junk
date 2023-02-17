package configuration

import "encoding/xml"

type XmlServerStorageTimeSettings struct {
	XMLName xml.Name `xml:"Time"`

	// Children.
	Format string `xml:"Format"`
	Zone   string `xml:"Zone"`
}
