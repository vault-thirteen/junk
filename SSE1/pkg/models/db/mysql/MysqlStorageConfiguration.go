package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/configuration"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/db/common"
)

// Configuration of a MySQL Storage.
type MysqlStorageConfiguration struct {

	// Common Settings.
	Address              string
	ConnectionParameters map[string]string
	Database             string
	Password             string
	User                 string

	// Data Source Name.
	// Is filled later, at Storage Initialization Stage.
	Dsn string

	// Initialization SQL Scripts' Settings and Sources.
	InitializationScripts configuration.ServerStorageIniScripts

	// Database Table Settings.
	TableSettings []common.TableSettings

	// Time Settings.
	Time configuration.ServerStorageTimeConfiguration

	// Other Settings taken from Application's Settings...

	// Timeout Interval Settings for Spam-Requests.
	CoolDownPeriods configuration.ServerStorageCoolDownPeriods

	// Idle Session Timeout Interval Setting.
	IdleSessionTimeoutSec uint

	// JSON Web Token Lifetime Duration Setting.
	TokenLifeTimeSec uint
}

// Configuration Constructor.
func NewMysqlStorageConfiguration(
	ssCfg configuration.ServerStorageConfiguration,
) (cfg *MysqlStorageConfiguration) {
	cfg = &MysqlStorageConfiguration{
		Address:              ssCfg.CommonParameters.Address,
		ConnectionParameters: ssCfg.CommonParameters.ConnectionParameters,
		Database:             ssCfg.CommonParameters.Database,
		Password:             ssCfg.CommonParameters.Password,
		User:                 ssCfg.CommonParameters.User,
		//
		Time:                  ssCfg.Time,
		TableSettings:         ssCfg.TableSettings,
		InitializationScripts: ssCfg.InitializationScripts,
		//
		CoolDownPeriods:       ssCfg.CoolDownPeriods,
		IdleSessionTimeoutSec: ssCfg.IdleSessionTimeoutSec,
		TokenLifeTimeSec:      ssCfg.TokenLifeTimeSec,
	}
	cfg.setDsn()
	return
}

// Sets the DSN of a Storage using its Settings.
func (c *MysqlStorageConfiguration) setDsn() {
	var cfg mysql.Config
	cfg = mysql.Config{
		Net:    MysqlStorageNetDefault,
		Addr:   c.Address,
		Params: c.ConnectionParameters,
		DBName: c.Database,
		User:   c.User,
		Passwd: c.Password,
		Loc:    c.Time.Zone,
	}
	c.Dsn = cfg.FormatDSN()
}
