package configuration

import "encoding/xml"

type XmlServerStorageTableSettings struct {
	XMLName xml.Name `xml:"TableSettings"`

	// Children.
	Table []XmlServerStorageTableSettingsTable `xml:"Table"`
}
