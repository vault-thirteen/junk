package database

import (
	"database/sql"
	"errors"
	"spdb_manager/database/query"
	"spdb_manager/models"
	"sync"

	"go.uber.org/multierr"
)

// DbForFilesImport is a helper-object for file import operations with database.
type DbForFilesImport struct {
	lock sync.Mutex
	db   *sql.DB
	tx   *sql.Tx

	// Prepared Statements.
	getServerIdStmt        *sql.Stmt
	getContainerTypeIdStmt *sql.Stmt
	getCodecIdStmt         *sql.Stmt
	insertNewMovieStmt     *sql.Stmt
}

// prepareStatements prepares SQL statements for the database.
func (dbo *DbForFilesImport) prepareStatements() (err error) {
	dbo.getServerIdStmt, err = dbo.tx.Prepare(query.GetServerIdByAddress)
	if err != nil {
		return err
	}

	dbo.getContainerTypeIdStmt, err = dbo.tx.Prepare(query.GetContainerTypeIdByName)
	if err != nil {
		return err
	}

	dbo.getCodecIdStmt, err = dbo.tx.Prepare(query.GetCodecIdByName)
	if err != nil {
		return err
	}

	dbo.insertNewMovieStmt, err = dbo.tx.Prepare(query.InsertNewMovie)
	if err != nil {
		return err
	}

	return nil
}

// Close disconnects from the database normally.
func (dbo *DbForFilesImport) Close() (err error) {
	return dbo.close(true)
}

// Abort disconnects from the database in case of emergency, i.e. on errors.
func (dbo *DbForFilesImport) Abort() (err error) {
	return dbo.close(false)
}

// A common method to disconnect from the database.
func (dbo *DbForFilesImport) close(isEverythingFine bool) (err error) {
	dbo.lock.Lock()
	defer dbo.lock.Unlock()

	defer func() {
		derr := dbo.db.Close()
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	defer func() {
		var derr error
		if isEverythingFine {
			derr = dbo.tx.Commit()
		} else {
			derr = dbo.tx.Rollback()
		}
		if derr != nil {
			err = multierr.Combine(err, derr)
		}
	}()

	// "The statements prepared for a transaction by calling the transaction's
	// Prepare or Stmt methods are closed by the call to Commit or Rollback."
	// Source: https://pkg.go.dev/database/sql#Tx.
	//
	// So, after we finish the Tx, we do not need to manually close all the
	// statements prepared within the Tx.

	return nil
}

// GetServerId fetches the server ID by its address.
func (dbo *DbForFilesImport) GetServerId(serverAddress string) (id uint, err error) {
	err = dbo.getServerIdStmt.QueryRow(serverAddress).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetContainerTypeId fetches the container type ID by its name.
func (dbo *DbForFilesImport) GetContainerTypeId(containerTypeName string) (id uint, err error) {
	err = dbo.getContainerTypeIdStmt.QueryRow(containerTypeName).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetCodecId fetches the codec ID by its name.
func (dbo *DbForFilesImport) GetCodecId(codecName string) (id uint, err error) {
	err = dbo.getCodecIdStmt.QueryRow(codecName).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// InsertNewMovie inserts a new movie into the database.
func (dbo *DbForFilesImport) InsertNewMovie(movie *models.MovieRaw) (err error) {
	var result sql.Result
	result, err = dbo.insertNewMovieStmt.Exec(
		movie.Title,
		movie.FileServerId,
		movie.FilePath,
		movie.FileName,
		movie.FileExtension,
		movie.ContainerTypeId,
		movie.OverallBitRate,
		movie.VideoCodecId,
		movie.AudioCodecId,
		movie.FrameWidth,
		movie.FrameHeight,
		movie.Duration,
	)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("no rows are affected")
	}

	return nil
}
