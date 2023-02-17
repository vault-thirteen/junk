package configuration

import "encoding/xml"

type XmlServerStorageTableSettingsTableColumn struct {
	XMLName xml.Name `xml:"Column"`

	// Attributes.
	Name string `xml:"name,attr"`
}
