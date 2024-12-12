package dbo

import (
	"database/sql"
	"errors"
	"fmt"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/models/settings"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"
	m "github.com/vault-thirteen/auxie/SQL/mysql"
)

const (
	SqlScriptFileExt                     = "sql"
	FileExtensionSeparator               = "."
	TableNamePrefixSeparator             = "_"
	SearchPattern_CreateTableIfNotExists = "CREATE TABLE IF NOT EXISTS %s"
	LastInsertedIdOnError                = -1
	CountOnError                         = -1
	IdOnError                            = 0
)

const (
	ErrTableNameIsNotFound = "table name is not found"
	ErrFRowsAffectedCount  = "affected rows count error: %v vs %v"
)

// DatabaseObject is a common database object, an entity managing all
// operations related to database.
type DatabaseObject struct {
	// Database settings.
	settings cs.DbSettings

	// SQL database object.
	db *sql.DB

	// Access control structure.
	guard *sync.RWMutex

	// List of prepared SQL statements.
	preparedStatements []*sql.Stmt

	// Source code of prepared statements.
	preparedStatementQueries []string
}

func NewDatabaseObject(settings cs.DbSettings) (dbo *DatabaseObject) {
	return &DatabaseObject{
		settings:                 settings,
		db:                       new(sql.DB),
		guard:                    new(sync.RWMutex),
		preparedStatements:       make([]*sql.Stmt, 0),
		preparedStatementQueries: make([]string, 0),
	}
}

// Init connects to the database, initialises the tables and prepares SQL
// statements. Before calling this method, a user must manually initialise
// table names and a list of query strings which will be used as prepared
// statements.
func (dbo *DatabaseObject) Init(preparedStatementQueryStrings []string) (err error) {
	fmt.Print(c.MsgConnectingToDatabase)

	err = dbo.connect()
	if err != nil {
		return err
	}

	err = dbo.initTables()
	if err != nil {
		return err
	}

	err = dbo.prepareStatements(preparedStatementQueryStrings)
	if err != nil {
		return err
	}

	fmt.Println(c.MsgOK)

	return nil
}

func (dbo *DatabaseObject) prepareStatements(preparedStatementQueryStrings []string) (err error) {
	dbo.preparedStatementQueries = make([]string, 0, len(preparedStatementQueryStrings))
	for _, pqs := range preparedStatementQueryStrings {
		dbo.preparedStatementQueries = append(dbo.preparedStatementQueries, pqs)
	}

	dbo.preparedStatements = make([]*sql.Stmt, 0, len(dbo.preparedStatementQueries))
	var st *sql.Stmt
	for _, psq := range dbo.preparedStatementQueries {
		st, err = dbo.db.Prepare(psq)
		if err != nil {
			return err
		}

		dbo.preparedStatements = append(dbo.preparedStatements, st)
	}

	return nil
}

func (dbo *DatabaseObject) PrefixTableName(tableName string) (tableNameFull string) {
	if len(dbo.settings.TableNamePrefix) > 0 {
		return dbo.settings.TableNamePrefix + TableNamePrefixSeparator + tableName
	}

	return tableName
}

func (dbo *DatabaseObject) connect() (err error) {
	mc := mysql.Config{
		Net:                  dbo.settings.Net,
		Addr:                 net.JoinHostPort(dbo.settings.Host, strconv.FormatUint(uint64(dbo.settings.Port), 10)),
		DBName:               dbo.settings.DBName,
		User:                 dbo.settings.User,
		Passwd:               dbo.settings.Password,
		AllowNativePasswords: dbo.settings.AllowNativePasswords,
		CheckConnLiveness:    dbo.settings.CheckConnLiveness,
		MaxAllowedPacket:     dbo.settings.MaxAllowedPacket,
		Params:               dbo.settings.Params,
	}

	dbo.db, err = sql.Open(dbo.settings.DriverName, mc.FormatDSN())
	if err != nil {
		return err
	}

	err = dbo.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (dbo *DatabaseObject) initTables() (err error) {
	var tableName string
	var tableExists bool

	// Create those tables which are not found.
	for _, tableNameOriginal := range dbo.settings.TablesToInit {
		tableName = dbo.PrefixTableName(tableNameOriginal)

		tableExists, err = m.TableExists(dbo.db, dbo.settings.DBName, tableName)
		if err != nil {
			return err
		}

		if !tableExists {
			log.Println(fmt.Sprintf(c.MsgFTableIsNotFound, tableName))

			err = dbo.initTable(tableName, tableNameOriginal)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// initTable runs initialisation scripts for tables which require
// initialisation as per configuration.
// tableName is a prefixed table name.
// tableNameOriginal is an original table name.
func (dbo *DatabaseObject) initTable(tableName string, tableNameOriginal string) (err error) {
	log.Println(fmt.Sprintf(c.MsgFInitialisingDatabaseTable, tableName))

	sqlScriptFilePath := filepath.Join(dbo.settings.TableInitScriptsFolder, tableNameOriginal+FileExtensionSeparator+SqlScriptFileExt)

	var buf []byte
	buf, err = os.ReadFile(sqlScriptFilePath)
	if err != nil {
		return err
	}

	var cmd string
	cmd, err = dbo.replaceTableNameInCreateTableScript(buf, tableName, tableNameOriginal)
	if err != nil {
		return err
	}

	_, err = dbo.db.Exec(cmd)
	if err != nil {
		return err
	}

	return nil
}

// replaceTableNameInCreateTableScript replaces a table name in the SQL script
// which creates a table.
// scriptText is contents of an SQL script file.
// tableName is a prefixed table name.
// tableNameOriginal is an original table name.
func (dbo *DatabaseObject) replaceTableNameInCreateTableScript(scriptText []byte, tableName string, tableNameOriginal string) (cmd string, err error) {
	pattern := fmt.Sprintf(SearchPattern_CreateTableIfNotExists, tableNameOriginal)
	replacement := fmt.Sprintf(SearchPattern_CreateTableIfNotExists, tableName)
	scriptTextStr := string(scriptText)

	if strings.Index(scriptTextStr, pattern) < 0 {
		return "", errors.New(ErrTableNameIsNotFound)
	}

	return strings.Replace(scriptTextStr, pattern, replacement, 1), nil
}

// ProbeDb pings the database server.
func (dbo *DatabaseObject) ProbeDb() (err error) {
	return dbo.db.Ping()
}

// Fin disconnects from the database.
func (dbo *DatabaseObject) Fin() (err error) {
	return dbo.db.Close()
}

// LockForReading locks database access for reading operations.
func (dbo *DatabaseObject) LockForReading() {
	dbo.guard.RLock()
}

// UnlockAfterReading unlocks database access after reading operations.
func (dbo *DatabaseObject) UnlockAfterReading() {
	dbo.guard.RUnlock()
}

// LockForWriting locks database access for writing operations.
func (dbo *DatabaseObject) LockForWriting() {
	dbo.guard.Lock()
}

// UnlockAfterWriting unlocks database access after writing operations.
func (dbo *DatabaseObject) UnlockAfterWriting() {
	dbo.guard.Unlock()
}

func (dbo *DatabaseObject) PreparedStatement(psIndex int) *sql.Stmt {
	return dbo.preparedStatements[psIndex]
}

func (dbo *DatabaseObject) DB() *sql.DB {
	return dbo.db
}
