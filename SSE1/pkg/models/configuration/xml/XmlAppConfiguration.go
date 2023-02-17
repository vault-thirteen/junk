package configuration

import (
	"encoding/xml"
	"strings"

	"github.com/vault-thirteen/junk/SSE1/pkg/helper/file"
)

type XmlAppConfiguration struct {
	XMLName xml.Name `xml:"Configuration"`

	// Children.
	Server XmlServerConfiguration `xml:"Server"`
}

func NewXmlAppConfiguration(
	filePath string,
) (xmlCfg *XmlAppConfiguration, err error) {

	// Get File Contents.
	var cfgFileContents string
	cfgFileContents, err = file.GetTextFileContents(filePath)
	if err != nil {
		return
	}

	// Decode the File.
	var decoder = xml.NewDecoder(strings.NewReader(cfgFileContents))
	decoder.Strict = true
	xmlCfg = new(XmlAppConfiguration)
	err = decoder.Decode(xmlCfg)
	if err != nil {
		return
	}
	return
}
