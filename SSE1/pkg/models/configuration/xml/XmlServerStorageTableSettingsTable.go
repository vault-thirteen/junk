package configuration

import "encoding/xml"

type XmlServerStorageTableSettingsTable struct {
	XMLName xml.Name `xml:"Table"`

	// Attributes.
	Name string `xml:"name,attr"`

	// Children.
	Column []XmlServerStorageTableSettingsTableColumn `xml:"Column"`
}
