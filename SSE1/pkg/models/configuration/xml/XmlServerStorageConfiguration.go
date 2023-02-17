package configuration

import "encoding/xml"

type XmlServerStorageConfiguration struct {
	XMLName xml.Name `xml:"Storage"`

	// Attributes.
	Address              string `xml:"address,attr"`
	ConnectionParameters string `xml:"connectionParameters,attr"`
	Database             string `xml:"database,attr"`
	Password             string `xml:"password,attr"`
	Type                 string `xml:"type,attr"`
	User                 string `xml:"user,attr"`

	// Children.
	InitializationScripts XmlServerStorageIniScripts    `xml:"InitializationScripts"`
	TableSettings         XmlServerStorageTableSettings `xml:"TableSettings"`
	Time                  XmlServerStorageTimeSettings  `xml:"Time"`
}
