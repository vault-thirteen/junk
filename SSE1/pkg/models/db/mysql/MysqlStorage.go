package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kr/pretty"
	vmysql "github.com/vault-thirteen/SQL/mysql"
	vtime "github.com/vault-thirteen/auxie/time"
	"github.com/vault-thirteen/errorz"
	jwtHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/jwt"
	loggerHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/logger"
	"github.com/vault-thirteen/junk/SSE1/pkg/interfaces/logger"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/configuration"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/http/request"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/session"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/user"
)

// Settings.
const (
	MysqlStorageNetDefault           = "tcp"
	MysqlDriver                      = "mysql"
	MysqlCommonTransactionTimeoutSec = 60
)

// Errors.
const (
	ErrCanNotInsertBrowserUserAgent  = "Can not insert a Browser User Agent"
	ErrConnectionIsNull              = "Connection is Null"
	ErrCoolDown                      = "You need to cool down"
	ErrCanNotLogInWithActiveSessions = "Can not log in with active Session(s) running"
	ErrCanNotCreateSession           = "Can not create a Session"
	ErrUserNotFound                  = "User is not found"
	ErrNoSessionToUpdate             = "No Session to update"
)

// Error Formats.
const (
	ErrfCanNotDisableUser        = "Can not disable a User with Id %v"
	ErrfTableColumnNamesMismatch = "Table Column Names Mismatch: %v"
	ErrfTableDoesNotExist        = "Table '%v' does not exist"
	ErrfTimeOffsetMismatch       = "Storage Time Offset Mismatch: %v vs %v"
)

// Messages.
const (
	MsgfInitializingTable = "Initializing a Table '%v'"
)

// Frequently used Query Templates.
const (
	QueryGetLastInsertedId = `SELECT LAST_INSERT_ID();`
)

// MySQL Storage.
// Implements the 'IStorage' Interface.
type MysqlStorage struct {
	Configuration MysqlStorageConfiguration
	db            *sql.DB
	logger        logger.ILogger
}

// Storage Constructor.
func NewMysqlStorage(
	cfg configuration.ServerStorageConfiguration,
	logger logger.ILogger,
) (storage *MysqlStorage, err error) {
	storage = &MysqlStorage{
		Configuration: *NewMysqlStorageConfiguration(cfg),
		logger:        logger,
	}
	return
}

// Reads the Data Source Name of the Storage.
func (s *MysqlStorage) GetDsn() string {
	return s.Configuration.Dsn
}

// Connects the Storage.
func (s *MysqlStorage) Connect() (err error) {
	s.db, err = sql.Open(MysqlDriver, s.Configuration.Dsn)
	if err != nil {
		return
	}
	err = s.db.Ping()
	if err != nil {
		return
	}
	return
}

// Checks the Integrity Storage.
func (s *MysqlStorage) Check() (err error) {
	err = s.db.Ping()
	if err != nil {
		return
	}
	err = s.checkTables()
	if err != nil {
		return
	}
	err = s.checkTime()
	if err != nil {
		return
	}
	return
}

// Disconnects the Storage.
func (s *MysqlStorage) Disconnect() (err error) {
	if s.db == nil {
		err = errors.New(ErrConnectionIsNull)
		return
	}
	err = s.db.Close()
	if err != nil {
		return
	}
	return
}

// Checks whether a User Name exists in the Database.
// The State of a User (registered or disabled) is ignored.
func (s *MysqlStorage) IsUserAuthenticationNameFree(
	userAuthenticationName string,
) (nameIsFree bool, err error) {
	const QueryCount = `SELECT COUNT(*)
FROM UserAuthentication
WHERE (Name = ?);`
	var statement *sql.Stmt
	statement, err = s.db.Prepare(QueryCount)
	if err != nil {
		return
	}
	defer func() {
		var derr = statement.Close()
		err = errorz.Combine(err, derr)
	}()
	var row = statement.QueryRow(userAuthenticationName)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return
	}
	if count == 0 {
		nameIsFree = true
	}
	return
}

// Check if a registered User with Id exists.
func (s *MysqlStorage) RegisteredUserIdExists(
	userId uint,
) (exists bool, err error) {
	const QueryCount = `SELECT COUNT(*)
FROM User
WHERE
	(Id = ?) AND
	(IsEnabled = true);`
	var statement *sql.Stmt
	statement, err = s.db.Prepare(QueryCount)
	if err != nil {
		return
	}
	defer func() {
		var derr = statement.Close()
		err = errorz.Combine(err, derr)
	}()
	var row = statement.QueryRow(userId)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return
	}
	if count == 1 {
		exists = true
	}
	return
}

// Reads a User's Id by its Authentication Name.
func (s *MysqlStorage) GetUserIdByAuthenticationName(
	userAuthenticationName string,
) (userId uint, err error) {
	const Query = `SELECT UserId
FROM UserAuthentication
WHERE Name = ?;`
	var row = s.db.QueryRow(Query, userAuthenticationName)
	err = row.Scan(&userId)
	if err != nil {
		return
	}
	return
}

// Registers a User.
func (s *MysqlStorage) RegisterUser(
	usr *user.User,
) (err error) {
	const (
		QueryfInsertUser = `INSERT INTO User 
(PublicName) 
VALUES (?);`
		QueryfInsertUserAuthentication = `INSERT INTO UserAuthentication 
(UserId, Name, Password) 
VALUES (?, ?, ?);`
		QueryfInsertUserRegistration = `INSERT INTO UserRegistration 
(UserId, SecretCode) 
VALUES (?, ?);`
	)

	// Start a Transaction.
	var tx *sql.Tx
	var cancelFunc context.CancelFunc
	tx, cancelFunc, err = s.startTx()
	if err != nil {
		return
	}
	defer cancelFunc()

	// Insert a User and get its Id.
	_, err = tx.Exec(QueryfInsertUser, usr.PublicName)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	var row = tx.QueryRow(QueryGetLastInsertedId)
	err = row.Scan(&usr.Id)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Insert a User Authentication.
	_, err = tx.Exec(
		QueryfInsertUserAuthentication,
		usr.Id,
		usr.Authentication.Name,
		usr.Authentication.Password,
	)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Insert a User Registration.
	_, err = tx.Exec(
		QueryfInsertUserRegistration,
		usr.Id,
		usr.Registration.SecretCode,
	)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Commit a Transaction.
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

// Disables a User.
// When a User is disabled, it can not be registered again.
// To register a User with the disabled Name,
// the disabled User must be first deleted from the Database.
func (s *MysqlStorage) DisableUser(
	usr *user.User,
) (err error) {
	var timePassedSinceLUAT int
	timePassedSinceLUAT, err = s.updateUserRegistrationLastUnregAttemptTimeByUserId(usr.Id)
	if err != nil {
		return
	}
	if timePassedSinceLUAT < int(s.Configuration.CoolDownPeriods.UserUnreg) {
		err = errors.New(ErrCoolDown)
		return
	}
	err = s.disableUser(usr)
	if err != nil {
		return
	}
	return
}

// Logs a User in.
func (s *MysqlStorage) LogUserIn(
	usr *user.User,
	machine *request.UserLogRequestMachine,
) (ses *session.Session, token *jwtHelper.TokenData, err error) {
	var timePassedSinceLLIAT int
	timePassedSinceLLIAT, err = s.updateUserAuthenticationLastLogInAttemptTimeByUserId(usr.Id)
	if err != nil {
		return
	}
	if timePassedSinceLLIAT < int(s.Configuration.CoolDownPeriods.UserLogIn) {
		err = errors.New(ErrCoolDown)
		return
	}
	ses, token, err = s.logUserIn(usr, machine)
	if err != nil {
		return
	}
	return
}

// Logs a User out.
func (s *MysqlStorage) LogUserOut(
	usr *user.User,
	ses *session.Session,
) (err error) {

	const (
		QueryFinishSession = `UPDATE Session
SET 
	EndTime=now()
WHERE
	(Id = ?) AND
	(UserId = ?) AND
  	(UserHost = ?) AND
	(UserBuaId = ?) AND
	(StartTime = ?) AND
	(EndTime IS NULL) AND
    (Marker = ?) AND
	(MarkerHash = ?) AND
	(TokenKey = ?);`
	)

	// Start a Transaction.
	var tx *sql.Tx
	var cancelFunc context.CancelFunc
	tx, cancelFunc, err = s.startTx()
	if err != nil {
		return
	}
	defer cancelFunc()

	_, err = tx.Exec(
		QueryFinishSession,
		ses.Id,
		usr.Id,
		ses.User.Host,
		ses.User.BrowserUserAgentId,
		ses.StartTime,
		ses.Marker,
		ses.MarkerHash,
		ses.TokenKey,
	)
	if err != nil {
		return
	}

	// Commit a Transaction.
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

// Lists public Names of all registered Users.
func (s *MysqlStorage) ListRegisteredUsersPublicNames() (pubNames []string, err error) {
	const (
		QuerySelect = `SELECT PublicName
FROM User
WHERE
	(IsEnabled = true);`
	)

	var rows *sql.Rows
	rows, err = s.db.Query(QuerySelect)
	if err != nil {
		return
	}
	pubNames = make([]string, 0)
	var pubName string
	for rows.Next() {
		err = rows.Scan(&pubName)
		if err != nil {
			return
		}
		pubNames = append(pubNames, pubName)
	}
	return
}

// Returns an Id of a Browser's User Agent by its full Name.
func (s *MysqlStorage) GetBrowserUserAgentId(
	browserUserAgentName string,
) (id uint, err error) {

	const (
		QueryCount = `SELECT count(*)
FROM BrowserUserAgent
WHERE Name = ?;`
		QuerySelect = `SELECT Id
FROM BrowserUserAgent
WHERE Name = ?;`
		QueryInsert = `INSERT INTO BrowserUserAgent (Name) 
VALUES (?);`
	)

	// Start a Transaction.
	var tx *sql.Tx
	var cancelFunc context.CancelFunc
	tx, cancelFunc, err = s.startTx()
	if err != nil {
		return
	}
	defer cancelFunc()

	// Search for a Record.
	var row *sql.Row
	var recordsCount int
	row = tx.QueryRow(QueryCount, browserUserAgentName)
	err = row.Scan(&recordsCount)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	if recordsCount == 1 {
		// Get the existing Record.
		row = tx.QueryRow(QuerySelect, browserUserAgentName)
		err = row.Scan(&id)
		if err != nil {
			err = rollbackTx(tx, err)
			return
		}
		// Commit a Transaction.
		err = tx.Commit()
		if err != nil {
			return
		}
		return
	}

	// Add the Record and return its Id.
	var sqlResult sql.Result
	sqlResult, err = tx.Exec(QueryInsert, browserUserAgentName)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	var rowsAffected int64
	rowsAffected, err = sqlResult.RowsAffected()
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	if rowsAffected != 1 {
		err = errors.New(ErrCanNotInsertBrowserUserAgent)
		err = rollbackTx(tx, err)
		return
	}
	row = tx.QueryRow(QueryGetLastInsertedId)
	err = row.Scan(&id)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	// Commit a Transaction.
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

// Returns a Key for a Token with a specified Marker Hash Sum.
// Token's Session must be active.
func (s *MysqlStorage) GetTokenKeyByMarkerHash(
	markerHash string,
) (tokenKey interface{}, err error) {

	const (
		QuerySelect = `SELECT TokenKey
FROM Session
WHERE
	(MarkerHash = ?) AND
	(EndTime IS NULL) AND
	(TIMESTAMPDIFF(SECOND, LastAccessTime, now()) < ?);`
	)

	var row = s.db.QueryRow(
		QuerySelect,
		markerHash,
		s.Configuration.TokenLifeTimeSec,
	)
	var tokenKeyStr string
	err = row.Scan(&tokenKeyStr)
	if err != nil {
		return
	}
	tokenKey = tokenKeyStr
	return
}

// Returns an active Session by its Id.
func (s *MysqlStorage) GetActiveSessionById(
	sessionId uint,
) (ses *session.Session, err error) {

	const (
		QuerySelect = `SELECT
	Id,
	UserId,
	UserHost,
	UserBuaId,
	StartTime,
	LastAccessTime,
	EndTime,
	Marker,
	MarkerHash,
	TokenKey
FROM Session
WHERE
	(Id = ?) AND
	(EndTime IS NULL) AND
	(TIMESTAMPDIFF(SECOND, LastAccessTime, now()) < ?);`
	)

	var row = s.db.QueryRow(
		QuerySelect,
		sessionId,
		s.Configuration.IdleSessionTimeoutSec,
	)
	var startTimeStr string
	var lastAccessTime string
	var endTimeStr sql.NullString
	ses = new(session.Session)
	err = row.Scan(
		&ses.Id,
		&ses.User.Id,
		&ses.User.Host,
		&ses.User.BrowserUserAgentId,
		&startTimeStr,
		&lastAccessTime,
		&endTimeStr,
		&ses.Marker,
		&ses.MarkerHash,
		&ses.TokenKey,
	)
	if err != nil {
		return
	}
	ses.StartTime, err = s.parseStorageTimeString(startTimeStr)
	if err != nil {
		return
	}
	ses.LastAccessTime, err = s.parseStorageTimeString(lastAccessTime)
	if err != nil {
		return
	}
	ses.EndTime, err = s.parseStorageTimeNullString(endTimeStr)
	if err != nil {
		return
	}
	return
}

// Updates the 'LastAccessTime' Parameter of an active Session.
func (s *MysqlStorage) UpdateActiveSessionLastAccessTime(
	ses *session.Session,
) (err error) {

	const (
		QueryUpdate = `UPDATE Session
SET LastAccessTime=now()
WHERE 
	(Id = ?) AND
	(UserId = ?) AND
	(UserHost = ?) AND
	(UserBuaId = ?) AND
	(StartTime = ?) AND
	(Marker = ?) AND
	(MarkerHash = ?) AND
	(TokenKey = ?) AND
	(EndTime IS NULL) AND
	(TIMESTAMPDIFF(SECOND, LastAccessTime, now()) < ?);`
	)

	// Start a Transaction.
	var tx *sql.Tx
	var cancelFunc context.CancelFunc
	tx, cancelFunc, err = s.startTx()
	if err != nil {
		return
	}
	defer cancelFunc()

	// Update the Time of the last Access.
	var sqlResult sql.Result
	sqlResult, err = tx.Exec(
		QueryUpdate,
		ses.Id,
		ses.User.Id,
		ses.User.Host,
		ses.User.BrowserUserAgentId,
		ses.StartTime,
		ses.Marker,
		ses.MarkerHash,
		ses.TokenKey,
		s.Configuration.TokenLifeTimeSec,
	)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	var rowsAffected int64
	rowsAffected, err = sqlResult.RowsAffected()
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	if rowsAffected != 1 {
		err = errors.New(ErrNoSessionToUpdate)
		err = rollbackTx(tx, err)
		return
	}

	// Commit a Transaction.
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) disableUser(
	usr *user.User,
) (err error) {

	const (
		QueryCount = `SELECT count(*)
FROM UserAuthentication AS ua
JOIN UserRegistration AS ur ON ur.UserId = ua.UserId
JOIN User AS u ON u.Id = ua.UserId
WHERE
	(ua.UserId = ?) AND
    (u.IsEnabled = true) AND
	(ua.Name = ?) AND
    (ua.Password = ?) AND
	(ur.SecretCode= ?);`
		QueryUpdateUser = `UPDATE User
SET IsEnabled = 0
WHERE Id = ?;`
		QueryUpdateUserRegistration = `UPDATE UserRegistration
SET UnregTime = now()
WHERE UserId = ?;`
		QueryFinishSessions = `UPDATE Session
SET EndTime = now()
WHERE UserId = ?;`
	)

	var sqlResult sql.Result
	var rowsAffected int64

	// Start a Transaction.
	var tx *sql.Tx
	var cancelFunc context.CancelFunc
	tx, cancelFunc, err = s.startTx()
	if err != nil {
		return
	}
	defer cancelFunc()

	// Check Data.
	var row *sql.Row
	row = tx.QueryRow(
		QueryCount,
		usr.Id,
		usr.Authentication.Name,
		usr.Authentication.Password,
		usr.Registration.SecretCode,
	)
	var userCount int
	err = row.Scan(&userCount)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	if userCount != 1 {
		err = errors.New(ErrUserNotFound)
		err = rollbackTx(tx, err)
		return
	}

	// Disable a User.
	sqlResult, err = tx.Exec(QueryUpdateUser, usr.Id)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	rowsAffected, err = sqlResult.RowsAffected()
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	if rowsAffected != 1 {
		err = fmt.Errorf(ErrfCanNotDisableUser, usr.Id)
		err = rollbackTx(tx, err)
		return
	}

	// Set the Time of Disabling.
	sqlResult, err = tx.Exec(QueryUpdateUserRegistration, usr.Id)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	rowsAffected, err = sqlResult.RowsAffected()
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	if rowsAffected != 1 {
		err = fmt.Errorf(ErrfCanNotDisableUser, usr.Id)
		err = rollbackTx(tx, err)
		return
	}

	// Finish all the Sessions of a disabled User.
	_, err = tx.Exec(QueryFinishSessions, usr.Id)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Commit a Transaction.
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) checkTables() (err error) {
	for _, tableSetting := range s.Configuration.TableSettings {
		err = s.checkTable(tableSetting.TableName, tableSetting.TableColumnNames)
		if err != nil {
			return
		}
	}
	return
}

func (s *MysqlStorage) checkTime() (err error) {
	var timeOffsetInConfig int
	timeOffsetInConfig, err = vtime.GetLocationOffsetHours(s.Configuration.Time.Zone)
	if err != nil {
		return
	}
	var timeOffsetInStorage int
	timeOffsetInStorage, err = s.getTimeOffset()
	if err != nil {
		return
	}
	if timeOffsetInConfig != timeOffsetInStorage {
		err = fmt.Errorf(
			ErrfTimeOffsetMismatch,
			timeOffsetInConfig,
			timeOffsetInStorage,
		)
		return
	}
	return
}

func (s *MysqlStorage) getTimeOffset() (offsetHrs int, err error) {
	const Query = `SELECT TIMESTAMPDIFF(HOUR, UTC_TIMESTAMP, NOW());`
	var row = s.db.QueryRow(Query)
	err = row.Scan(&offsetHrs)
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) checkTable(
	tableName string,
	tableColumnNames []string,
) (err error) {
	err = s.checkTableExistence(tableName)
	if err != nil {
		// Try to Initialize the Table using the Script.
		var err2 = s.applyTableIniScript(tableName)
		if err2 != nil {
			err = errorz.Combine(err, err2)
			return
		} else {
			err = nil
		}
	}
	err = s.checkTableColumnsExistence(tableName, tableColumnNames)
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) applyTableIniScript(
	tableName string,
) (err error) {
	var message = fmt.Sprintf(MsgfInitializingTable, tableName)
	loggerHelper.UseLogger(s.logger, message)

	// Start a Transaction.
	var tx *sql.Tx
	var cancelFunc context.CancelFunc
	tx, cancelFunc, err = s.startTx()
	if err != nil {
		return
	}
	defer cancelFunc()

	// Perform a Script.
	var script = s.Configuration.InitializationScripts.TableScripts[tableName]
	_, err = tx.Exec(script)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Commit a Transaction.
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) checkTableExistence(
	tableName string,
) (err error) {
	var tableExists bool
	tableExists, err = vmysql.TableExists(
		s.db,
		s.Configuration.Database,
		tableName,
	)
	if err != nil {
		return
	}
	if !tableExists {
		err = fmt.Errorf(ErrfTableDoesNotExist, tableName)
		return
	}
	return
}

func (s *MysqlStorage) checkTableColumnsExistence(
	tableName string,
	tableColumnNames []string,
) (err error) {
	var tableColumnNamesFound []string
	tableColumnNamesFound, err = vmysql.GetTableColumnNames(
		s.db,
		s.Configuration.Database,
		tableName,
	)
	if err != nil {
		return
	}
	sort.Strings(tableColumnNames)
	var ok = reflect.DeepEqual(tableColumnNames, tableColumnNamesFound)
	if !ok {
		err = fmt.Errorf(
			ErrfTableColumnNamesMismatch,
			pretty.Diff(tableColumnNames, tableColumnNamesFound),
		)
		return
	}
	return
}

func (s *MysqlStorage) startTx() (tx *sql.Tx, cancelFunc context.CancelFunc, err error) {
	var ctx context.Context
	ctx, cancelFunc = context.WithTimeout(
		context.Background(),
		time.Second*MysqlCommonTransactionTimeoutSec,
	)
	tx, err = s.db.BeginTx(
		ctx,
		&sql.TxOptions{Isolation: sql.LevelSerializable},
	)
	if err != nil {
		return
	}
	return
}

func rollbackTx(
	tx *sql.Tx,
	errIn error,
) (errOut error) {
	var errRollback = tx.Rollback()
	errOut = errorz.Combine(errIn, errRollback)
	return
}

func (s *MysqlStorage) updateUserRegistrationLastUnregAttemptTimeByUserId(
	userId uint,
) (timePassedSinceLUAT int, err error) {
	const (
		QuerySelect = `SELECT TIMESTAMPDIFF(SECOND, COALESCE(LastUnregAttemptTime, RegTime), now())
FROM UserRegistration
WHERE UserId = ?;`
		QueryUpdate = `UPDATE UserRegistration
SET LastUnregAttemptTime=now()
WHERE UserId = ?;`
	)

	// Start a Transaction.
	var tx *sql.Tx
	var cancelFunc context.CancelFunc
	tx, cancelFunc, err = s.startTx()
	if err != nil {
		return
	}
	defer cancelFunc()

	// Calculate Time passed since the last Attempt to Disable a User.
	var row = tx.QueryRow(QuerySelect, userId)
	err = row.Scan(&timePassedSinceLUAT)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Update the Time of the last Disabling Attempt.
	_, err = tx.Exec(QueryUpdate, userId)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Commit a Transaction.
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) parseStorageTimeString(
	timeText string,
) (t time.Time, err error) {
	if len(timeText) == 0 {
		return
	}
	t, err = time.ParseInLocation(
		s.Configuration.Time.Format,
		timeText,
		s.Configuration.Time.Zone,
	)
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) parseStorageTimeNullString(
	timeText sql.NullString,
) (t sql.NullTime, err error) {
	if !timeText.Valid {
		t = sql.NullTime{
			Valid: false,
		}
		return
	}
	t = sql.NullTime{
		Valid: true,
	}
	t.Time, err = s.parseStorageTimeString(timeText.String)
	if err != nil {
		t.Valid = false
		return
	}
	return
}

func (s *MysqlStorage) updateUserAuthenticationLastLogInAttemptTimeByUserId(
	userId uint,
) (timePassedSinceLLIAT int, err error) {

	const (
		QuerySelect = `SELECT TIMESTAMPDIFF(SECOND, COALESCE(LastLogInAttemptTime, '2000-01-01'), now())
FROM UserAuthentication AS ua
JOIN UserRegistration AS ur ON ua.UserId = ur.UserId
WHERE ua.UserId = ?;`
		QueryUpdate = `UPDATE UserAuthentication
SET LastLogInAttemptTime=now()
WHERE UserId = ?;`
	)

	// Start a Transaction.
	var tx *sql.Tx
	var cancelFunc context.CancelFunc
	tx, cancelFunc, err = s.startTx()
	if err != nil {
		return
	}
	defer cancelFunc()

	// Calculate Time passed since the last Attempt to Log a User in.
	var row = tx.QueryRow(QuerySelect, userId)
	err = row.Scan(&timePassedSinceLLIAT)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Update the Time of the last Log-In Attempt.
	_, err = tx.Exec(QueryUpdate, userId)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Commit a Transaction.
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) logUserIn(
	usr *user.User,
	machine *request.UserLogRequestMachine,
) (ses *session.Session, token *jwtHelper.TokenData, err error) {

	const (
		QueryCountUsers = `SELECT count(*)
FROM UserAuthentication AS ua
JOIN UserRegistration AS ur ON ur.UserId = ua.UserId
JOIN User AS u ON u.Id = ua.UserId
WHERE
	(ua.UserId = ?) AND
    (u.IsEnabled = true) AND
	(ua.Name = ?) AND
    (ua.Password = ?);`
		QueryCountActiveSessions = `SELECT count(*)
FROM Session
WHERE
	(UserId = ?) AND
	(EndTime IS NULL) AND
	(TIMESTAMPDIFF(SECOND, LastAccessTime, now()) < ?);`
	)

	// Start a Transaction.
	var tx *sql.Tx
	var cancelFunc context.CancelFunc
	tx, cancelFunc, err = s.startTx()
	if err != nil {
		return
	}
	defer cancelFunc()

	// Check the User.
	var row *sql.Row
	row = tx.QueryRow(
		QueryCountUsers,
		usr.Id,
		usr.Authentication.Name,
		usr.Authentication.Password,
	)
	var userCount int
	err = row.Scan(&userCount)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	if userCount != 1 {
		err = errors.New(ErrUserNotFound)
		err = rollbackTx(tx, err)
		return
	}

	// Finish all the idle Sessions of a User.
	err = s.finishUserIdleSessions(tx, usr)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Check the active Sessions.
	row = tx.QueryRow(
		QueryCountActiveSessions,
		usr.Id,
		s.Configuration.IdleSessionTimeoutSec,
	)
	var activeSessionsCount int
	err = row.Scan(&activeSessionsCount)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}
	if activeSessionsCount != 0 {
		err = errors.New(ErrCanNotLogInWithActiveSessions)
		err = rollbackTx(tx, err)
		return
	}

	// Create the Session.
	ses, token, err = s.createSession(tx, usr, machine)
	if err != nil {
		err = rollbackTx(tx, err)
		return
	}

	// Commit a Transaction.
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) finishUserIdleSessions(
	tx *sql.Tx,
	usr *user.User,
) (err error) {
	const QueryFinishIdleSessions = `UPDATE Session
SET EndTime=now()
WHERE
	(UserId = ?) AND
	(EndTime IS NULL) AND
    (TIMESTAMPDIFF(SECOND, LastAccessTime, now()) > ?);`

	_, err = tx.Exec(
		QueryFinishIdleSessions,
		usr.Id,
		s.Configuration.IdleSessionTimeoutSec,
	)
	if err != nil {
		return
	}
	return
}

func (s *MysqlStorage) createSession(
	tx *sql.Tx,
	usr *user.User,
	machine *request.UserLogRequestMachine,
) (ses *session.Session, td *jwtHelper.TokenData, err error) {

	const (
		QueryCreateSession = `INSERT INTO Session (
	UserId, 
	UserHost, 
	UserBuaId, 
	StartTime,
	Marker,
	MarkerHash,
	TokenKey
) VALUES (?, ?, ?, ?, ?, ?, ?);`
	)

	// Prepare Data.
	// While Token stores an Id of a Session inside,
	// we can not create a Token before a Session.
	td, err = jwtHelper.PrepareDataForToken()
	if err != nil {
		return
	}
	ses = &session.Session{
		Marker:     td.UniqueMarker,
		MarkerHash: td.UniqueMarkerHash,
		TokenKey:   td.TokenKeyStr,
		//Id:         0, // Is taken from the Database.
		StartTime: time.Now(),
		User: session.SessionUser{
			Id:                 usr.Id,
			Host:               machine.Host,
			BrowserUserAgent:   machine.BrowserUserAgent.Name,
			BrowserUserAgentId: machine.BrowserUserAgent.Id,
		},
	}

	// Create a Session.
	var sqlResult sql.Result
	sqlResult, err = tx.Exec(
		QueryCreateSession,
		usr.Id,
		machine.Host,
		machine.BrowserUserAgent.Id,
		ses.StartTime,
		ses.Marker,
		ses.MarkerHash,
		ses.TokenKey,
	)
	if err != nil {
		return
	}
	var rowsAffected int64
	rowsAffected, err = sqlResult.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffected != 1 {
		err = errors.New(ErrCanNotCreateSession)
		return
	}
	var row *sql.Row
	row = tx.QueryRow(QueryGetLastInsertedId)
	err = row.Scan(&ses.Id)
	if err != nil {
		return
	}

	// Create a Token.
	var token *jwt.Token
	token, err = jwtHelper.CreateJWToken(
		ses.Id,
		ses.MarkerHash,
		s.Configuration.TokenLifeTimeSec,
	)
	if err != nil {
		return
	}
	td.TokenString, err = token.SignedString(td.TokenKey)
	if err != nil {
		return
	}
	return
}
