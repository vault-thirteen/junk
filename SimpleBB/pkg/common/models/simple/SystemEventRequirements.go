package simple

// SystemEventRequirements store requirements for system event data.
// Different event types require different sets of parameters.
type SystemEventRequirements struct {
	IsThreadIdRequired  bool
	IsMessageIdRequired bool
	IsUserIdRequired    bool
	IsCreatorRequired   bool
}
