package cm

import (
	"fmt"

	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
)

type ConfigurationService struct {
	ShortName  string
	FullName   string
	Components []ccse.CommonConfigurationServiceEntry
	Servers    []ccse.CommonConfigurationServiceEntry
	Clients    []ccse.CommonConfigurationServiceEntry
}

func getConfigurationServiceEntriesByType(entriesIn []ccse.CommonConfigurationServiceEntry, entryType string) (entriesOut []ccse.CommonConfigurationServiceEntry) {
	return filterConfigurationServiceEntries(entriesIn, Field_Type, entryType)
}
func getConfigurationServiceEntriesByProtocol(entriesIn []ccse.CommonConfigurationServiceEntry, entryProtocol string) (entriesOut []ccse.CommonConfigurationServiceEntry) {
	return filterConfigurationServiceEntries(entriesIn, Field_Protocol, entryProtocol)
}
func getConfigurationServiceEntryByTypeAndProtocol(entriesIn []ccse.CommonConfigurationServiceEntry, entryType string, entryProtocol string) (entry *ccse.CommonConfigurationServiceEntry) {
	es1 := getConfigurationServiceEntriesByType(entriesIn, entryType)
	es2 := getConfigurationServiceEntriesByProtocol(es1, entryProtocol)

	if len(es2) == 0 {
		return nil
	}

	return &es2[0]
}

func parseConfigurationServiceEntries(rawEntries []CommonConfigurationServiceEntryRaw) (es []ccse.CommonConfigurationServiceEntry, err error) {
	es = []ccse.CommonConfigurationServiceEntry{}
	var e *ccse.CommonConfigurationServiceEntry
	var p *ccp.CommonConfigurationParameter
	for _, rawEntry := range rawEntries {
		e = &ccse.CommonConfigurationServiceEntry{
			Type:       rawEntry.Type,
			Protocol:   rawEntry.Protocol,
			Parameters: []ccp.CommonConfigurationParameter{},
		}

		for _, rep := range rawEntry.Parameters {
			p = &ccp.CommonConfigurationParameter{
				Name: rep.Name,
				Type: rep.Type,
			}

			p.Value, err = ccp.ParseCommonConfigurationParameterValue(rep.Type, rep.Value)
			if err != nil {
				return nil, err
			}

			e.Parameters = append(e.Parameters, *p)
		}

		es = append(es, *e)
	}

	return es, nil
}

// filterConfigurationServiceEntries filters entries by a field.
// Supported fields for filtering are: type, protocol.
func filterConfigurationServiceEntries(entriesIn []ccse.CommonConfigurationServiceEntry, filterType string, filterValue string) (entriesOut []ccse.CommonConfigurationServiceEntry) {
	entriesOut = []ccse.CommonConfigurationServiceEntry{}

	switch filterType {
	case Field_Type:
		for _, e := range entriesIn {
			if e.Type == filterValue {
				entriesOut = append(entriesOut, e)
			}
		}
		return entriesOut

	case Field_Protocol:
		for _, e := range entriesIn {
			if e.Protocol == filterValue {
				entriesOut = append(entriesOut, e)
			}
		}
		return entriesOut
	}

	// This can not happen while this function is used internally.
	err := fmt.Errorf(ErrF_FilterTypeIsUnsupported, filterType)
	panic(err)
	return nil
}
