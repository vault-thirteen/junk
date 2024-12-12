package mode

const (
	WorkModeCreateDb         = "create"
	WorkModeInitDb           = "initialize"
	WorkModeImportVideosToDb = "import"
	WorkModeEditDb           = "edit"
	WorkModeViewDb           = "view"
)

var modes = map[string]bool{
	WorkModeCreateDb:         true,
	WorkModeInitDb:           true,
	WorkModeImportVideosToDb: true,
	WorkModeEditDb:           true,
	WorkModeViewDb:           true,
}

// IsAvailable checks whether the specified work mode is allowed to be used.
func IsAvailable(mode string) (ok bool) {
	_, ok = modes[mode]

	// Check existence.
	if !ok {
		return false
	}

	// Check state.
	return modes[mode]
}
