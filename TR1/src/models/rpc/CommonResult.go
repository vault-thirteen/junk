package rm

type CommonResult struct {
	// Time taken to perform the request, in milliseconds.
	TimeSpent int64 `json:"timeSpent,omitempty"`
}

// Clear sets values to zeroed state so that it is removed from JSON output.
func (cr CommonResult) Clear() {
	cr.TimeSpent = 0
}
