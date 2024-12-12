package dbo

const (
	TableIncidents     = "Incidents"
	TableNotifications = "Notifications"
	TableResources     = "Resources"
	TableSystemEvents  = "SystemEvents"
)

type TableNames struct {
	Incidents     string
	Notifications string
	Resources     string
	SystemEvents  string
}
