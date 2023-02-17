package settings

import (
	"encoding/xml"
	"strings"

	"github.com/vault-thirteen/junk/SSE2/internal/helper"
)

type XmlSettings struct {
	XMLName xml.Name `xml:"Settings"`

	FileSizeLimiter XmlSettingsFileSizeLimiter `xml:"FileSizeLimiter"`
}

type XmlSettingsFileSizeLimiter struct {
	XMLName xml.Name `xml:"FileSizeLimiter"`

	MimeType []XmlSettingsFileSizeLimiterMimeType `xml:"MimeType"`
}

type XmlSettingsFileSizeLimiterMimeType struct {
	XMLName xml.Name `xml:"MimeType"`

	Name      string `xml:"name,attr"`
	SizeLimit int    `xml:"sizeLimit,attr"`
}

func NewXmlSettings(
	filePath string,
) (xmlConfig *XmlSettings, err error) {
	var cfgFileContents string
	cfgFileContents, err = helper.GetTextFileContents(filePath)
	if err != nil {
		return nil, err
	}

	var decoder = xml.NewDecoder(strings.NewReader(cfgFileContents))
	decoder.Strict = true

	xmlConfig = new(XmlSettings)

	err = decoder.Decode(xmlConfig)
	if err != nil {
		return nil, err
	}

	return xmlConfig, nil
}
