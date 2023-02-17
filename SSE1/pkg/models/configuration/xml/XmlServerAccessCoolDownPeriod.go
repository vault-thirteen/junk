package configuration

import "encoding/xml"

type XmlServerAccessCoolDownPeriod struct {
	XMLName xml.Name `xml:"CoolDownPeriod"`

	// Children.
	UserLogInSec uint `xml:"userLogInSec,attr"`
	UserUnregSec uint `xml:"userUnregSec,attr"`
}
