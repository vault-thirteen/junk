package cm

import (
	"encoding/json"
)

type ConfigurationRaw struct {
	Service ConfigurationServiceRaw `json:"service"`
}

type ConfigurationServiceRaw struct {
	ShortName  string                               `json:"shortName"`
	FullName   string                               `json:"fullName"`
	Components []CommonConfigurationServiceEntryRaw `json:"components"`
	Servers    []CommonConfigurationServiceEntryRaw `json:"servers"`
	Clients    []CommonConfigurationServiceEntryRaw `json:"clients"`
}

type CommonConfigurationServiceEntryRaw struct {
	Type       string                            `json:"type"`
	Protocol   string                            `json:"protocol"`
	Parameters []CommonConfigurationParameterRaw `json:"parameters"`
}

type CommonConfigurationParameterRaw struct {
	Name  string          `json:"name"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}
