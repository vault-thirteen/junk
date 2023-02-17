// test.go.

package eom

import (
	"database/sql"

	// PostgreSQL Driver.
	_ "github.com/lib/pq"

	"github.com/vault-thirteen/SQL/postgresql"
)

// Test Database Parameters.
const (
	TestDatabaseDriver     = "postgres"
	TestDatabaseHost       = "localhost"
	TestDatabasePort       = "5432"
	TestDatabaseDatabase   = "test"
	TestDatabaseUser       = "test"
	TestDatabasePassword   = "test"
	TestDatabaseParameters = "sslmode=disable"
)

func makeTestDatabaseDsn() (dsn string) {
	dsn = postgresql.MakeDsn(
		TestDatabaseHost,
		TestDatabasePort,
		TestDatabaseDatabase,
		TestDatabaseUser,
		TestDatabasePassword,
		TestDatabaseParameters,
	)
	return
}

func connectToTestDatabase(
	dsn string,
) (sqlConnection *sql.DB, err error) {
	return sql.Open(TestDatabaseDriver, dsn)
}
