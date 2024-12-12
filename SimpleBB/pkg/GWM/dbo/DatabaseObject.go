package dbo

import (
	as "github.com/vault-thirteen/SimpleBB/pkg/ACM/settings"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/models/dbo"
)

type DatabaseObject struct {
	cdbo.DatabaseObject

	// List of prefixed table names.
	tableNames *TableNames
}

func NewDatabaseObject(settings as.DbSettings) (dbo *DatabaseObject) {
	commonDBO := cdbo.NewDatabaseObject(settings)

	dbo = &DatabaseObject{}
	dbo.DatabaseObject = *commonDBO
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
		IPAddresses: dbo.prefixTableName(TableIPAddresses),
	}
}

func (dbo *DatabaseObject) prefixTableName(tableName string) (tableNameFull string) {
	return dbo.DatabaseObject.PrefixTableName(tableName)
}
