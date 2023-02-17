// ElectionOfMaster_test.go.

package eom

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/vault-thirteen/SQL/postgresql"
	"github.com/vault-thirteen/tester"
)

func Test_NewElectionOfMaster(t *testing.T) {

	const QueryDropTableTest = `DROP TABLE "test_election";`

	var aTest *tester.Test
	var dsn string
	var eom *ElectionOfMaster
	var err error
	var sqlConnection *sql.DB

	aTest = tester.New(t)

	dsn = makeTestDatabaseDsn()
	sqlConnection, err = connectToTestDatabase(dsn)
	eom, err = NewElectionOfMaster(
		"test",
		"first",
		ElectionTypeSingleMaster,
		dsn,
		60,
	)
	aTest.MustBeNoError(err)
	aTest.MustBeDifferent(eom, nil)

	// Check the internal Fields.
	aTest.MustBeEqual(eom.serviceName, "test")
	aTest.MustBeEqual(eom.serviceInstanceID, "first")
	aTest.MustBeEqual(eom.kind, ElectionType(ElectionTypeSingleMaster))
	aTest.MustBeEqual(eom.sqlConnectionDsn, dsn)
	aTest.MustBeEqual(eom.tableName, "test_election")
	aTest.MustBeEqual(eom.procedureName, "election_of_master")
	aTest.MustBeEqual(eom.expirationPeriodSec, uint(60))

	// Finalization.
	// Check the SQL Connection.
	sqlConnection, err = sql.Open(DatabaseDriver, dsn)
	aTest.MustBeNoError(err)
	_, err = sqlConnection.Exec(QueryDropTableTest)
	aTest.MustBeNoError(err)
	err = sqlConnection.Close()
	aTest.MustBeNoError(err)
}

func Test_ElectMyself(t *testing.T) {

	const QueryDropTableTest = `DROP TABLE "myservice_election";`

	var aTest *tester.Test
	var dsn string
	var electionResult ElectionFeedback
	var electionResultTime time.Time
	var eom_1 *ElectionOfMaster
	var eom_2 *ElectionOfMaster
	var err error
	var ok bool
	var serviceName string
	var sqlConnection *sql.DB

	aTest = tester.New(t)

	dsn = makeTestDatabaseDsn()
	serviceName = "myservice"

	// Prepare two virtual elected Service Instances (S.I.) with different IDs.
	eom_1, err = NewElectionOfMaster(
		serviceName,
		"first",
		ElectionTypeSingleMaster,
		dsn,
		10,
	)
	aTest.MustBeNoError(err)
	eom_2, err = NewElectionOfMaster(
		serviceName,
		"second",
		ElectionTypeSingleMaster,
		dsn,
		10,
	)
	aTest.MustBeNoError(err)

	// Do the Test...

	// 1. Take the Lease by the first S.I.
	electionResult, err = eom_1.ElectMyself()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(electionResult.Result, true)
	electionResultTime = electionResult.Time

	// 2. Try to take the Lease by the second S.I.
	time.Sleep(time.Second * 3)
	electionResult, err = eom_2.ElectMyself()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(electionResult.Result, false)
	aTest.MustBeEqual(electionResult.Time, electionResultTime)

	// 3. Try to renew the Lease by the first S.I.
	time.Sleep(time.Second * 2)
	electionResult, err = eom_1.ElectMyself()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(electionResult.Result, true)
	ok = (electionResult.Time.Sub(electionResultTime) >= time.Second*5) &&
		(electionResult.Time.Sub(electionResultTime) < time.Second*6)
	aTest.MustBeEqual(ok, true)
	electionResultTime = electionResult.Time

	// 4. Wait for the Lease to expire and take it by the second S.I.
	time.Sleep(time.Second * (10 + 1))
	electionResult, err = eom_2.ElectMyself()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(electionResult.Result, true)
	ok = (electionResult.Time.Sub(electionResultTime) >= time.Second*11) &&
		(electionResult.Time.Sub(electionResultTime) < time.Second*12)
	aTest.MustBeEqual(ok, true)

	// Finalization.
	time.Sleep(time.Second * 1)
	sqlConnection, err = sql.Open("postgres", dsn)
	aTest.MustBeNoError(err)
	_, err = sqlConnection.Exec(QueryDropTableTest)
	aTest.MustBeNoError(err)
	err = sqlConnection.Close()
	aTest.MustBeNoError(err)
}

func Test_DestroyInfrastructure(t *testing.T) {

	var aTest *tester.Test
	var dsn string
	var eom *ElectionOfMaster
	var err error
	var procedureExists bool
	var sqlConnection *sql.DB
	var tableExists bool

	aTest = tester.New(t)

	dsn = makeTestDatabaseDsn()
	eom, err = NewElectionOfMaster(
		"test",
		"first",
		ElectionTypeSingleMaster,
		dsn,
		60,
	)
	aTest.MustBeNoError(err)
	aTest.MustBeDifferent(eom, nil)

	err = eom.DestroyInfrastructure()
	aTest.MustBeNoError(err)

	// Verification.
	sqlConnection, err = sql.Open("postgres", dsn)
	aTest.MustBeNoError(err)

	tableExists, err = postgresql.TableExists(sqlConnection, TableSchema, eom.tableName)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(tableExists, false)

	procedureExists, err = postgresql.ProcedureExists(sqlConnection, TableSchema, eom.procedureName)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(procedureExists, false)

	// Finalization.
	err = sqlConnection.Close()
	aTest.MustBeNoError(err)
}
