package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/vault-thirteen/SQL/postgresql"
	"go.uber.org/multierr"
)

// Database Settings.
const (
	DatabaseDriver     = "postgres"
	DatabaseHost       = "localhost"
	DatabasePort       = "5432"
	DatabaseUser       = "test"
	DatabasePassword   = "test"
	DatabaseDatabase   = "test"
	DatabaseParameters = "sslmode=disable"
)

const (
	ErrRowHasNotBeenInserted = "row has not been inserted"
	ErrItemsCountMismatch    = "items count mismatch"
)

const QueryFInsertItem = `INSERT INTO %v ("Name", "IsSpecial")
VALUES ($1, $2);`
const QueryFReadSpecialItems = `SELECT * FROM %v WHERE "IsSpecial" = true;`
const DataSize = 5500
const SpecialItemRarity = 1000 // One in 1000.
const ItemsCountExpected = 5

func main() {
	log.Printf(
		"Data Size: %v, Special Item Rarity: 1 of %v.\r\n",
		DataSize,
		SpecialItemRarity,
	)

	var err error
	err = fillTablesWithData(DataSize)
	mustBeNoError(err)

	log.Println("Sleeping ...")
	time.Sleep(time.Second * 10)

	var workTimes []time.Duration
	workTimes, err = readSpecialItemsFromTables()
	mustBeNoError(err)

	for i, workTime := range workTimes {
		log.Printf("Work time #%v: %v mcs.\r\n", i+1, workTime.Microseconds())
	}
}

func fillTablesWithData(dataSize int) (err error) {
	var db *sql.DB
	db, err = connect()
	if err != nil {
		return err
	}

	defer func() {
		derr := disconnect(db)
		err = multierr.Combine(err, derr)
	}()

	log.Println("Filling the A table ...")
	err = fillTable(db, "boolean_index_table_a", dataSize)
	if err != nil {
		return err
	}

	log.Println("Filling the B table ...")
	err = fillTable(db, "boolean_index_table_b", dataSize)
	if err != nil {
		return err
	}

	return nil
}

func fillTable(db *sql.DB, tableName string, dataSize int) (err error) {
	queryTemplate := fmt.Sprintf(QueryFInsertItem, tableName)

	var statement *sql.Stmt
	statement, err = db.Prepare(queryTemplate)
	if err != nil {
		return err
	}

	defer func() {
		derr := statement.Close()
		err = multierr.Combine(err, derr)
	}()

	// Insert rows.
	var (
		result        sql.Result
		rowsAffected  int64
		itemName      string
		itemIsSpecial bool
	)

	for i := 1; i <= dataSize; i++ {
		itemName = fmt.Sprintf("name_%v", i)
		if i%SpecialItemRarity != 0 {
			itemIsSpecial = false
		} else {
			itemIsSpecial = true
		}

		result, err = statement.Exec(itemName, itemIsSpecial)
		if err != nil {
			return err
		}

		rowsAffected, err = result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected != 1 {
			return errors.New(ErrRowHasNotBeenInserted)
		}
	}

	return nil
}

func readSpecialItemsFromTables() (workTimes []time.Duration, err error) {
	var db *sql.DB
	db, err = connect()
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := disconnect(db)
		err = multierr.Combine(err, derr)
	}()

	workTimes = make([]time.Duration, 2)

	workTimes[0], err = readTable(db, "boolean_index_table_a")
	if err != nil {
		return nil, err
	}

	workTimes[1], err = readTable(db, "boolean_index_table_b")
	if err != nil {
		return nil, err
	}

	return workTimes, nil
}

func readTable(db *sql.DB, tableName string) (workTime time.Duration, err error) {
	timeStart := time.Now()
	defer func() {
		workTime = time.Since(timeStart)
	}()

	queryTemplate := fmt.Sprintf(QueryFReadSpecialItems, tableName)

	var statement *sql.Stmt
	statement, err = db.Prepare(queryTemplate)
	if err != nil {
		return workTime, err
	}

	defer func() {
		derr := statement.Close()
		err = multierr.Combine(err, derr)
	}()

	var rows *sql.Rows
	rows, err = statement.Query()
	if err != nil {
		return workTime, err
	}

	// Simulate reading the items.
	itemsCount := 0
	for rows.Next() {
		itemsCount++
	}

	if itemsCount != ItemsCountExpected {
		return workTime, errors.New(ErrItemsCountMismatch)
	}

	return workTime, nil
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func connect() (db *sql.DB, err error) {
	dsn := postgresql.MakeDsn(
		DatabaseHost,
		DatabasePort,
		DatabaseDatabase,
		DatabaseUser,
		DatabasePassword,
		DatabaseParameters,
	)

	db, err = sql.Open(DatabaseDriver, dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func disconnect(db *sql.DB) (err error) {
	return db.Close()
}
