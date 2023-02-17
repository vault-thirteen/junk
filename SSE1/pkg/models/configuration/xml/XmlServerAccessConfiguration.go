package configuration

import "encoding/xml"

type XmlServerAccessConfiguration struct {
	XMLName xml.Name `xml:"Access"`

	// Children.
	CoolDownPeriod XmlServerAccessCoolDownPeriod       `xml:"CoolDownPeriod"`
	Session        XmlServerAccessSessionConfiguration `xml:"Session"`
	Token          XmlServerAccessTokenConfiguration   `xml:"Token"`
}
