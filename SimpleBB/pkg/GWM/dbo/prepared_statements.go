package dbo

import (
	"database/sql"
	"fmt"
)

// Indices of prepared statements.
const (
	DbPsid_CountBlocksByIPAddress = 0
	DbPsid_AddBlock               = 1
	DbPsid_IncreaseBlockDuration  = 2
	DbPsid_ClearIPAddresses       = 3
)

func (dbo *DatabaseObject) makePreparedStatementQueryStrings() (qs []string) {
	var q string
	qs = make([]string, 0)

	// 0.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE UserIPAB = ?;`, dbo.tableNames.IPAddresses)
	qs = append(qs, q)

	// 1.
	q = fmt.Sprintf(`INSERT INTO %s (UserIPAB, TIME) VALUES (?, DATE_ADD(Now(), INTERVAL ? SECOND));`, dbo.tableNames.IPAddresses)
	qs = append(qs, q)

	// 2.
	q = fmt.Sprintf(`UPDATE %s SET TIME = DATE_ADD(TIME, INTERVAL ? SECOND) WHERE UserIPAB = ?;`, dbo.tableNames.IPAddresses)
	qs = append(qs, q)

	// 3.
	q = fmt.Sprintf(`DELETE FROM %s WHERE TIME < Now();`, dbo.tableNames.IPAddresses)
	qs = append(qs, q)

	return qs
}

func (dbo *DatabaseObject) GetPreparedStatementByIndex(i int) (ps *sql.Stmt) {
	return dbo.DatabaseObject.PreparedStatement(i)
}
