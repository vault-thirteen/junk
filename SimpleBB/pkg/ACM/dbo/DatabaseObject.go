package dbo

import (
	as "github.com/vault-thirteen/SimpleBB/pkg/ACM/settings"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/models/dbo"
)

type DatabaseObject struct {
	cdbo.DatabaseObject

	// System parameters.
	sp SystemParameters

	// List of prefixed table names.
	tableNames *TableNames
}

func NewDatabaseObject(settings as.DbSettings, sp SystemParameters) (dbo *DatabaseObject) {
	commonDBO := cdbo.NewDatabaseObject(settings)

	dbo = &DatabaseObject{}
	dbo.DatabaseObject = *commonDBO
	dbo.sp = sp
	dbo.tableNames = new(TableNames)

	return dbo
}

// Init connects to the database, initialises the tables and prepares SQL
// statements.
func (dbo *DatabaseObject) Init() (err error) {
	dbo.initTableNames()

	var preparedStatementQueryStrings = dbo.makePreparedStatementQueryStrings()

	err = dbo.DatabaseObject.Init(preparedStatementQueryStrings)
	if err != nil {
		return err
	}

	return nil
}

func (dbo *DatabaseObject) initTableNames() {
	dbo.tableNames = &TableNames{
		PreRegisteredUsers: dbo.prefixTableName(TablePreRegisteredUsers),
		Users:              dbo.prefixTableName(TableUsers),
		PreSessions:        dbo.prefixTableName(TablePreSessions),
		Sessions:           dbo.prefixTableName(TableSessions),
		Incidents:          dbo.prefixTableName(TableIncidents),
		PasswordChanges:    dbo.prefixTableName(TablePasswordChanges),
		EmailChanges:       dbo.prefixTableName(TableEmailChanges),
		LogEvents:          dbo.prefixTableName(TableLogEvents),
	}
}

func (dbo *DatabaseObject) prefixTableName(tableName string) (tableNameFull string) {
	return dbo.DatabaseObject.PrefixTableName(tableName)
}
