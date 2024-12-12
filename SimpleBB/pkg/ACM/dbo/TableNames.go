package dbo

const (
	TablePreRegisteredUsers = "PreRegisteredUsers"
	TableUsers              = "Users"
	TablePreSessions        = "PreSessions"
	TableSessions           = "Sessions"
	TableIncidents          = "Incidents"
	TablePasswordChanges    = "PasswordChanges"
	TableEmailChanges       = "EmailChanges"
	TableLogEvents          = "LogEvents"
)

type TableNames struct {
	PreRegisteredUsers string
	Users              string
	PreSessions        string
	Sessions           string
	Incidents          string
	PasswordChanges    string
	EmailChanges       string
	LogEvents          string
}
