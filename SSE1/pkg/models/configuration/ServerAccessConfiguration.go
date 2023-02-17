package configuration

type ServerAccessConfiguration struct {
	CoolDownPeriod ServerAccessCoolDownPeriod
	Session        ServerAccessSessionConfiguration
	Token          ServerAccessTokenConfiguration
}
