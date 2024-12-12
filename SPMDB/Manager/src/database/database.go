package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// OpenSimple connects to a database for simple operations.
func OpenSimple(dbFile string) (dbo *DbSimple, err error) {
	dbo = new(DbSimple)

	dbo.db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	return dbo, nil
}

// OpenForFilesImport connects to a database for file import operations.
func OpenForFilesImport(dbFile string) (dbo *DbForFilesImport, err error) {
	dbo = new(DbForFilesImport)

	dbo.db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	dbo.tx, err = dbo.db.Begin()
	if err != nil {
		return nil, err
	}

	err = dbo.prepareStatements()
	if err != nil {
		return nil, err
	}

	return dbo, nil
}

// CheckConnection checks connection to a database by connecting to it and then
// disconnecting from it.
func CheckConnection(dbFile string) (err error) {
	var dbo *DbSimple
	dbo, err = OpenSimple(dbFile)
	if err != nil {
		return err
	}

	err = dbo.Close()
	if err != nil {
		return err
	}

	return nil
}
