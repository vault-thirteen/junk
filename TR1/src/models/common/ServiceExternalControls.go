package cm

type ServiceExternalControls struct {
	// Channel for an external controller. When a message comes from this
	// channel, an external controller must stop this server. Normally, the
	// server does not stop itself.
	mustBeStopped chan bool
}

func NewExternalControls() *ServiceExternalControls {
	return &ServiceExternalControls{
		mustBeStopped: make(chan bool),
	}
}
