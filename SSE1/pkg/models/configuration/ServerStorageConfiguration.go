package configuration

import (
	configuration "github.com/vault-thirteen/junk/SSE1/pkg/models/configuration/xml"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/db/common"
)

type ServerStorageConfiguration struct {
	Type                  ServerStorageType
	CommonParameters      common.StorageConfiguration
	InitializationScripts ServerStorageIniScripts
	TableSettings         []common.TableSettings
	Time                  ServerStorageTimeConfiguration

	// Settings taken from Application's Settings.
	CoolDownPeriods       ServerStorageCoolDownPeriods
	IdleSessionTimeoutSec uint
	TokenLifeTimeSec      uint
}

func NewTableSettings(
	settings configuration.XmlServerStorageTableSettings,
) (result []common.TableSettings, err error) {
	result = make([]common.TableSettings, 0, len(settings.Table))
	for _, table := range settings.Table {
		columnNames := make([]string, 0, len(table.Column))
		for _, column := range table.Column {
			columnNames = append(columnNames, column.Name)
		}
		result = append(result, common.TableSettings{
			TableName:        table.Name,
			TableColumnNames: columnNames,
		})
	}
	return
}
