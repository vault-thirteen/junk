package stn

// Errors.
const (
	ErrCLAFolderNotSet        = "folder is not set"
	ErrCLAFileWithNamesNotSet = "file with names is not set"
)

type Settings struct {
	FolderPath        string
	FileWithNamesPath string
}
