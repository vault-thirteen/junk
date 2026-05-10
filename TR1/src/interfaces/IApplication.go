package interfaces

type IApplication interface {
	GetConfiguration() IConfiguration
	Use() error
}
