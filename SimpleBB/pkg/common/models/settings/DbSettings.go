package s

import (
	"errors"
)

const (
	ErrDbDriverName             = "DB driver name is not set"
	ErrDbNet                    = "DB net is not set"
	ErrDbHost                   = "DB host is not set"
	ErrDbPort                   = "DB port is not set"
	ErrDbName                   = "DB name is not set"
	ErrDbUser                   = "DB user is not set"
	ErrDbMaxAllowedPacket       = "DB MaxAllowedPacket is not set"
	ErrDbTableInitScriptsFolder = "DB TableInitScriptsFolder is not set"
)

// DbSettings are parameters of a MySQL database together with some other
// settings related to database infrastructure. When a password is not set, it
// is taken from the stdin.
type DbSettings struct {
	// Access settings.
	DriverName string `json:"driverName"`
	Net        string `json:"net"`
	Host       string `json:"host"`
	Port       uint16 `json:"port"`
	DBName     string `json:"dbName"`
	User       string `json:"user"`
	Password   string `json:"password"`

	// Various specific MySQL settings.
	AllowNativePasswords bool              `json:"allowNativePasswords"`
	CheckConnLiveness    bool              `json:"checkConnLiveness"`
	MaxAllowedPacket     int               `json:"maxAllowedPacket"`
	Params               map[string]string `json:"params"`

	// Database structure and initialisation settings.
	TableNamePrefix        string   `json:"tableNamePrefix"`
	TablesToInit           []string `json:"tablesToInit"`
	TableInitScriptsFolder string   `json:"tableInitScriptsFolder"`
}

func (dbs DbSettings) Check() (err error) {
	if len(dbs.DriverName) == 0 {
		return errors.New(ErrDbDriverName)
	}
	if len(dbs.Net) == 0 {
		return errors.New(ErrDbNet)
	}
	if len(dbs.Host) == 0 {
		return errors.New(ErrDbHost)
	}
	if dbs.Port == 0 {
		return errors.New(ErrDbPort)
	}
	if len(dbs.DBName) == 0 {
		return errors.New(ErrDbName)
	}
	if len(dbs.User) == 0 {
		return errors.New(ErrDbUser)
	}
	if dbs.MaxAllowedPacket == 0 {
		return errors.New(ErrDbMaxAllowedPacket)
	}
	if len(dbs.TableInitScriptsFolder) == 0 {
		return errors.New(ErrDbTableInitScriptsFolder)
	}

	return nil
}
