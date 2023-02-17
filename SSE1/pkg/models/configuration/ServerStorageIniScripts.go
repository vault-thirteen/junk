package configuration

type ServerStorageIniScripts struct {

	// Path to a Folder containing a 'Table' Folder with Scripts.
	Folder string

	// SQL Scripts used for Initialization of Database Tables.
	TableScripts map[string]string
}
