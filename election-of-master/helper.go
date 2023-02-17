// helper.go.

package eom

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vault-thirteen/SQL/postgresql"
	"github.com/vault-thirteen/auxie/unicode"
	"github.com/vault-thirteen/errorz"
)

const (
	QueryFileExtFull      = ".sql"
	QueryFilesStoragePath = "sql/public/functions"
)

func checkServiceName(
	serviceName string,
) error {

	if len(serviceName) == 0 {
		return errors.New(ErrSizeIsNull)
	}

	for _, letter := range serviceName {
		if (!unicode.SymbolIsLatLetter(letter)) &&
			(!unicode.SymbolIsNumber(letter)) &&
			(letter != '_') {
			return fmt.Errorf(ErrfForbiddenSymbol, string(letter))
		}
	}

	return nil
}

func createElectionTable(
	connection *sql.DB,
	tableName string,
) (err error) {

	var query string
	var statement *sql.Stmt
	var tableNameIsGood bool

	// Verify the Table Name.
	tableNameIsGood, err = postgresql.TableNameIsGood(tableName)
	if (err != nil) || !tableNameIsGood {
		return err
	}

	// Prepare the Query using the dynamic SQL.
	query = fmt.Sprintf(QueryfCreateElectionTable, tableName)
	statement, err = connection.Prepare(query)
	if err != nil {
		return err
	}
	defer func() {
		var derr error
		derr = statement.Close()
		err = errorz.Combine(err, derr)
	}()

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

func deleteElectionTable(
	connection *sql.DB,
	tableName string,
) (err error) {

	var query string
	var statement *sql.Stmt
	var tableNameIsGood bool

	// Verify the Table Name.
	tableNameIsGood, err = postgresql.TableNameIsGood(tableName)
	if (err != nil) || !tableNameIsGood {
		return err
	}

	// Prepare the Query using the dynamic SQL.
	query = fmt.Sprintf(QueryfDeleteElectionTable, tableName)
	statement, err = connection.Prepare(query)
	if err != nil {
		return err
	}
	defer func() {
		var derr error
		derr = statement.Close()
		err = errorz.Combine(err, derr)
	}()

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

func createProcedure(
	connection *sql.DB,
	procedureName string,
) (err error) {

	var procedureNameIsGood bool
	var query string
	var queryPtr *string
	var queryFilePath string
	var statement *sql.Stmt

	// Verify the Procedure Name.
	procedureNameIsGood, err = postgresql.ProcedureNameIsGood(procedureName)
	if (err != nil) || !procedureNameIsGood {
		return err
	}

	// Get the Query Source Code.
	queryFilePath = filepath.Join(
		QueryFilesStoragePath,
		procedureName+QueryFileExtFull,
	)
	queryPtr, err = getFileContentsString(queryFilePath)
	if (err != nil) || (queryPtr == nil) {
		return err
	}
	query = *queryPtr

	// Prepare the Query using the dynamic SQL.
	statement, err = connection.Prepare(query)
	if err != nil {
		return err
	}
	defer func() {
		var derr error
		derr = statement.Close()
		err = errorz.Combine(err, derr)
	}()

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

func deleteProcedure(
	connection *sql.DB,
	procedureName string,
) (err error) {

	var query string
	var statement *sql.Stmt
	var procedureNameIsGood bool

	// Verify the Procedure Name.
	procedureNameIsGood, err = postgresql.ProcedureNameIsGood(procedureName)
	if (err != nil) || !procedureNameIsGood {
		return err
	}

	// Prepare the Query using the dynamic SQL.
	query = fmt.Sprintf(QueryfDeleteProcedure, procedureName)
	statement, err = connection.Prepare(query)
	if err != nil {
		return err
	}
	defer func() {
		var derr error
		derr = statement.Close()
		err = errorz.Combine(err, derr)
	}()

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

func getFileContentsString(
	filePath string,
) (contents *string, err error) {

	var bytes []byte
	var contentsString string
	var file *os.File

	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		var derr error
		derr = file.Close()
		err = errorz.Combine(err, derr)
	}()

	bytes, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	contentsString = string(bytes)
	return &contentsString, nil
}
