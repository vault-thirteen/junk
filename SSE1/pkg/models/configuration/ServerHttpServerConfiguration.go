package configuration

type ServerHttpServerConfiguration struct {
	Address            string
	CookiePath         string
	ShutdownTimeoutSec uint
	TokenHeader        string
}
