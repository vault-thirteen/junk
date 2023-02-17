package configuration

import "encoding/xml"

type XmlServerConfiguration struct {
	XMLName xml.Name `xml:"Server"`

	// Children.
	Access     XmlServerAccessConfiguration     `xml:"Access"`
	HttpServer XmlServerHttpServerConfiguration `xml:"HttpServer"`
	Logger     XmlServerLoggerConfiguration     `xml:"Logger"`
	Storage    XmlServerStorageConfiguration    `xml:"Storage"`
	TimeZone   string                           `xml:"TimeZone"`
	TLS        XmlServerTlsConfiguration        `xml:"TLS"`
}
