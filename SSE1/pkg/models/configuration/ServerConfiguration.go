package configuration

type ServerConfiguration struct {
	Access     ServerAccessConfiguration
	HttpServer ServerHttpServerConfiguration
	Logger     ServerLoggerConfiguration
	Storage    ServerStorageConfiguration
	Timezone   ServerTimezone
	TLS        ServerTlsConfiguration
}
