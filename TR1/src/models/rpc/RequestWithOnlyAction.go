package rm

import "encoding/json"

// RequestWithOnlyAction is a data set for an API request received from a
// client in JSON format. This data is later mixed with additional client data.
type RequestWithOnlyAction struct {
	// Action is a name of the function which a client wants to perform.
	Action *string `json:"action"`

	// Parameters are function parameters.
	Parameters *json.RawMessage `json:"parameters"`
}
