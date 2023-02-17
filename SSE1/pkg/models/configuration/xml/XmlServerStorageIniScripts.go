package configuration

import "encoding/xml"

type XmlServerStorageIniScripts struct {
	XMLName xml.Name `xml:"InitializationScripts"`

	// Attributes.
	Folder string `xml:"folder,attr"`
}
