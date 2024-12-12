package dbo

import (
	"database/sql"
	"fmt"
)

// Indices of prepared statements.
const (
	DbPsid_CountUserSubscriptions        = 0
	DbPsid_CountThreadSubscriptions      = 1
	DbPsid_InitUserSubscriptions         = 2
	DbPsid_InitThreadSubscriptions       = 3
	DbPsid_GetUserSubscriptions          = 4
	DbPsid_GetThreadSubscriptions        = 5
	DbPsid_SaveUserSubscriptions         = 6
	DbPsid_SaveThreadSubscriptions       = 7
	DbPsid_ClearThreadSubscriptionRecord = 8
	DbPsid_GetAllThreadSubscriptions     = 9
	DbPsid_GetAllUserSubscriptions       = 10
)

func (dbo *DatabaseObject) makePreparedStatementQueryStrings() (qs []string) {
	var q string
	qs = make([]string, 0)

	// 0.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE UserId = ?;`, dbo.tableNames.UserSubscriptions)
	qs = append(qs, q)

	// 1.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE ThreadId = ?;`, dbo.tableNames.ThreadSubscriptions)
	qs = append(qs, q)

	// 2.
	q = fmt.Sprintf(`INSERT INTO %s (UserId) VALUES (?);`, dbo.tableNames.UserSubscriptions)
	qs = append(qs, q)

	// 3.
	q = fmt.Sprintf(`INSERT INTO %s (ThreadId) VALUES (?);`, dbo.tableNames.ThreadSubscriptions)
	qs = append(qs, q)

	// 4.
	q = fmt.Sprintf(`SELECT Id, UserId, Threads FROM %s WHERE UserId = ?;`, dbo.tableNames.UserSubscriptions)
	qs = append(qs, q)

	// 5.
	q = fmt.Sprintf(`SELECT Id, ThreadId, Users FROM %s WHERE ThreadId = ?;`, dbo.tableNames.ThreadSubscriptions)
	qs = append(qs, q)

	// 6.
	q = fmt.Sprintf(`UPDATE %s SET Threads = ? WHERE UserId = ?;`, dbo.tableNames.UserSubscriptions)
	qs = append(qs, q)

	// 7.
	q = fmt.Sprintf(`UPDATE %s SET Users = ? WHERE ThreadId = ?;`, dbo.tableNames.ThreadSubscriptions)
	qs = append(qs, q)

	// 8.
	q = fmt.Sprintf(`DELETE FROM %s WHERE ThreadId = ?;`, dbo.tableNames.ThreadSubscriptions)
	qs = append(qs, q)

	// 9.
	q = fmt.Sprintf(`SELECT Id, ThreadId, Users FROM %s;`, dbo.tableNames.ThreadSubscriptions)
	qs = append(qs, q)

	// 10.
	q = fmt.Sprintf(`SELECT Id, UserId, Threads FROM %s;`, dbo.tableNames.UserSubscriptions)
	qs = append(qs, q)

	return qs
}

func (dbo *DatabaseObject) GetPreparedStatementByIndex(i int) (ps *sql.Stmt) {
	return dbo.DatabaseObject.PreparedStatement(i)
}
