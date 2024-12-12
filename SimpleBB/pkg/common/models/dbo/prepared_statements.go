package dbo

// Query format strings of prepared statements.
const (
	Query_SaveIncident               = `INSERT INTO %s (Module, Type, Email, UserIPAB) VALUES (?, ?, ?, ?);`
	Query_SaveIncidentWithoutUserIPA = `INSERT INTO %s (Module, Type, Email) VALUES (?, ?, ?);`
)
