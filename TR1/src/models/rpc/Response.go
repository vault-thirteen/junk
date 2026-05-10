package rm

type Response struct {
	// Action is a name of the function which was called to get this response.
	Action *string `json:"action"`

	// Result returned by the called function.
	Result any `json:"result"`
}
