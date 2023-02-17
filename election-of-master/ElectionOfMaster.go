// ElectionOfMaster.go.

package eom

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vault-thirteen/SQL/postgresql"
	"github.com/vault-thirteen/errorz"
)

// Election of Master.

// This Mechanism helps to have only a single Instance of the Service to be
// running at a Time.
//
// An existing Master must prove its Status by sending a Confirmation before
// the Expiration Period of the last Confirmation. If a Confirmation
// is expired, i.e. no one sends a Confirmation before the Expiration Period,
// then the first received Confirmation assigns the Master, i.e. the first One
// who sends the Confirmation, becomes the Master.

// Notes:
//
// 1.	Due to a Bug in the 'pq' Library, we have to open a new SQL Connection
// 		for each Query. This is a temporary Solution.

// Error Messages.
const (
	ErrfServiceName            = "Service Name '%v' is bad. Reason: %v."
	ErrSizeIsNull              = "size is Null"
	ErrfForbiddenSymbol        = "symbol '%v' is forbidden"
	ErrForbiddenElectionType   = "Forbidden Election Type"
	ErrSqlConnectionDsnIsEmpty = "SQL Connection DSN is not set"
)

const (
	TableSchema        = "public"
	TableNamePostfix   = "_election"
	TableProcedureName = "election_of_master"
	DatabaseDriver     = "postgres"

	QueryfCreateElectionTable = `CREATE TABLE "%s"
(
	"ServiceInstanceId" VARCHAR(255) UNIQUE NOT NULL,
	"LastUpdateTime" timestamp with time zone
);`
	QueryfDeleteElectionTable = `DROP TABLE "%s";`
	QueryfDeleteProcedure     = `DROP FUNCTION IF EXISTS %s;`
	QueryfElection            = `SELECT %v($1, $2, $3);`
)

type ElectionOfMaster struct {
	serviceName         string
	serviceInstanceID   string
	kind                ElectionType
	sqlConnectionDsn    string
	tableName           string
	procedureName       string
	expirationPeriodSec uint
}

// NewElectionOfMaster Function creates a new Election of a Master.
func NewElectionOfMaster(

	// Name of the Service. It affects the Name of a Table.
	serviceName string,

	// The unique ID of the Service Instance which takes Part in the Elections.
	serviceInstanceID string,

	kind ElectionType,
	sqlConnectionDsn string,
	expirationPeriodSec uint,
) (*ElectionOfMaster, error) {

	var eom *ElectionOfMaster
	var err error
	var procedureExists bool
	var sqlConnection *sql.DB
	var tableExists bool

	eom = new(ElectionOfMaster)

	// Connection DSN.
	if len(sqlConnectionDsn) == 0 {
		return nil, errors.New(ErrSqlConnectionDsnIsEmpty)
	}
	eom.sqlConnectionDsn = sqlConnectionDsn

	// Service Name.
	err = checkServiceName(serviceName)
	if err != nil {
		return nil, fmt.Errorf(ErrfServiceName, serviceName, err.Error())
	}
	eom.serviceName = serviceName

	// Service Instance ID.
	eom.serviceInstanceID = serviceInstanceID

	// Election Type.
	if !kind.IsValid() {
		return nil, errors.New(ErrForbiddenElectionType)
	}
	eom.kind = kind

	// Check the SQL Connection.
	sqlConnection, err = sql.Open(DatabaseDriver, eom.sqlConnectionDsn)
	if err != nil {
		return nil, err
	}
	err = sqlConnection.Ping()
	if err != nil {
		return nil, err
	}

	// Prepare the Database Table.
	eom.tableName = eom.serviceName + TableNamePostfix
	tableExists, err = postgresql.TableExists(sqlConnection, TableSchema, eom.tableName)
	if err != nil {
		return nil, err
	}
	if !tableExists {
		// Create the Table.
		err = createElectionTable(sqlConnection, eom.tableName)
		if err != nil {
			return nil, err
		}
	}

	// Prepare the Database Procedure.
	eom.procedureName = TableProcedureName
	procedureExists, err = postgresql.ProcedureExists(sqlConnection, TableSchema, eom.procedureName)
	if err != nil {
		return nil, err
	}
	if !procedureExists {
		// Create the Procedure.
		err = createProcedure(sqlConnection, eom.procedureName)
		if err != nil {
			return nil, err
		}
	}

	// Close the SQL Connection.
	err = sqlConnection.Close()
	if err != nil {
		return nil, err
	}

	// Various Settings.
	eom.expirationPeriodSec = expirationPeriodSec

	return eom, nil
}

// ElectMyself Method tries to take Part in the Election.
// Returns 'true' when we are elected as a Master, otherwise returns 'false'.
func (eom *ElectionOfMaster) ElectMyself() (result ElectionFeedback, err error) {

	var ctx context.Context
	var ctxCancelFunc context.CancelFunc
	var query string
	var row *sql.Row
	var sqlConnection *sql.DB
	var stmt *sql.Stmt
	var terr error
	var transaction *sql.Tx
	var transactionOptions sql.TxOptions

	// Check the SQL Connection.
	sqlConnection, err = sql.Open(DatabaseDriver, eom.sqlConnectionDsn)
	if err != nil {
		return result, err
	}
	err = sqlConnection.Ping()
	if err != nil {
		return result, err
	}

	// Prepare the Context.
	ctx, ctxCancelFunc = context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(eom.expirationPeriodSec)/time.Duration(2),
	)
	defer ctxCancelFunc()

	// Transaction Start.
	transactionOptions = sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	transaction, err = sqlConnection.BeginTx(ctx, &transactionOptions)
	if err != nil {
		return result, err
	}

	// Statement Preparation.
	query = fmt.Sprintf(QueryfElection, eom.procedureName)
	stmt, err = transaction.Prepare(query)
	if err != nil {
		return result, err
	}
	defer func() {
		var derr error
		derr = stmt.Close()
		err = errorz.Combine(err, derr)
	}()

	// Statement Execution.
	row = stmt.QueryRow(
		eom.tableName,
		eom.serviceInstanceID,
		eom.expirationPeriodSec,
	)
	err = row.Scan(&result.RawData)
	if err != nil {
		// Transaction Rollback.
		terr = transaction.Rollback()
		err = errorz.Combine(err, terr)
		return result, err
	}

	// Transaction Commit.
	err = transaction.Commit()
	if err != nil {
		return result, err
	}

	err = result.ParseRawData()
	if err != nil {
		return result, err
	}

	// Close the SQL Connection.
	err = sqlConnection.Close()
	if err != nil {
		return result, err
	}

	return result, err
}

func (eom *ElectionOfMaster) DestroyInfrastructure() (err error) {

	var sqlConnection *sql.DB

	// Check the SQL Connection.
	sqlConnection, err = sql.Open(DatabaseDriver, eom.sqlConnectionDsn)
	if err != nil {
		return err
	}
	err = sqlConnection.Ping()
	if err != nil {
		return err
	}

	// Delete the Procedure.
	err = deleteProcedure(sqlConnection, eom.procedureName)
	if err != nil {
		return err
	}

	// Delete the Table.
	err = deleteElectionTable(sqlConnection, eom.tableName)
	if err != nil {
		return err
	}

	// Close the SQL Connection.
	err = sqlConnection.Close()
	if err != nil {
		return err
	}

	return nil
}
