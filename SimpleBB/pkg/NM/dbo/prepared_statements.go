package dbo

import (
	"database/sql"
	"fmt"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/models/dbo"
)

// Indices of prepared statements.
const (
	DbPsid_InsertNewNotification            = 0
	DbPsid_GetNotificationById              = 1
	DbPsid_MarkNotificationAsRead           = 2
	DbPsid_DeleteNotificationById           = 3
	DbPsid_SaveIncident                     = 4
	DbPsid_SaveIncidentWithoutUserIPA       = 5
	DbPsid_GetAllNotificationsByUserId      = 6
	DbPsid_GetUnreadNotificationsByUserId   = 7
	DbPsid_CountUnreadNotificationsByUserId = 8
	DbPsid_ClearNotifications               = 9
	DbPsid_GetNotificationsByUserIdOnPage   = 10
	DbPsid_CountAllNotificationsByUserId    = 11
	DbPsid_SaveSystemEvent                  = 12
	DbPsid_GetSystemEventById               = 13
	DbPsid_AddResource                      = 14
	DbPsid_GetResourceById                  = 15
	DbPsid_DeleteResourceById               = 16
	DbPsid_ListAllResourceIdsOnPage         = 17
	DbPsid_CountAllResources                = 18
)

func (dbo *DatabaseObject) makePreparedStatementQueryStrings() (qs []string) {
	var q string
	qs = make([]string, 0)

	// 0.
	q = fmt.Sprintf(`INSERT INTO %s (UserId, Text, ToC) VALUES (?, ?, Now());`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 1.
	q = fmt.Sprintf(`SELECT Id, UserId, Text, ToC, IsRead, ToR FROM %s WHERE Id = ?;`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 2.
	q = fmt.Sprintf(`UPDATE %s SET IsRead = TRUE, ToR = Now() WHERE Id = ? AND UserId = ?;`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 3.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ? AND IsRead IS TRUE;`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 4.
	q = fmt.Sprintf(cdbo.Query_SaveIncident, dbo.tableNames.Incidents)
	qs = append(qs, q)

	// 5.
	q = fmt.Sprintf(cdbo.Query_SaveIncidentWithoutUserIPA, dbo.tableNames.Incidents)
	qs = append(qs, q)

	// 6.
	q = fmt.Sprintf(`SELECT Id, UserId, Text, ToC, IsRead, ToR FROM %s WHERE UserId = ?;`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 7.
	q = fmt.Sprintf(`SELECT Id, UserId, Text, ToC, IsRead, ToR FROM %s WHERE UserId = ? AND IsRead IS FALSE;`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 8.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE UserId = ? AND IsRead IS FALSE;`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 9.
	q = fmt.Sprintf(`DELETE FROM %s WHERE IsRead IS TRUE AND ToR < ?;`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 10.
	q = fmt.Sprintf(`SELECT Id, UserId, Text, ToC, IsRead, ToR FROM %s WHERE UserId = ? ORDER BY ToC DESC LIMIT ? OFFSET ?;`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 11.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE UserId = ?;`, dbo.tableNames.Notifications)
	qs = append(qs, q)

	// 12.
	q = fmt.Sprintf(`INSERT INTO %s (Type, ThreadId, MessageId, UserId) VALUES (?, ?, ?, ?);`, dbo.tableNames.SystemEvents)
	qs = append(qs, q)

	// 13.
	q = fmt.Sprintf(`SELECT Id, Type, ThreadId, MessageId, UserId, Time FROM %s WHERE Id = ?;`, dbo.tableNames.SystemEvents)
	qs = append(qs, q)

	// 14.
	q = fmt.Sprintf(`INSERT INTO %s (Type, Text, Number) VALUES (?, ?, ?);`, dbo.tableNames.Resources)
	qs = append(qs, q)

	// 15.
	q = fmt.Sprintf(`SELECT Id, Type, Text, Number, ToC FROM %s WHERE Id = ?;`, dbo.tableNames.Resources)
	qs = append(qs, q)

	// 16.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ?;`, dbo.tableNames.Resources)
	qs = append(qs, q)

	// 17.
	q = fmt.Sprintf(`SELECT Id FROM %s ORDER BY Id ASC LIMIT ? OFFSET ?;`, dbo.tableNames.Resources)
	qs = append(qs, q)

	// 18.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s;`, dbo.tableNames.Resources)
	qs = append(qs, q)

	return qs
}

func (dbo *DatabaseObject) GetPreparedStatementByIndex(i int) (ps *sql.Stmt) {
	return dbo.DatabaseObject.PreparedStatement(i)
}
