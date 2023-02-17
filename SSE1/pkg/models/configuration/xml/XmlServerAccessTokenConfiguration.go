package configuration

import "encoding/xml"

type XmlServerAccessTokenConfiguration struct {
	XMLName xml.Name `xml:"Token"`

	// Children.
	LifeTimeSec uint `xml:"lifeTimeSec,attr"`
}
