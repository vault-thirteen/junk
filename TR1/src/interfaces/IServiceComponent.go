package interfaces

type IServiceComponent interface {
	Init(cfg IConfiguration, controller IController) (sc IServiceComponent, err error)
	GetConfiguration() IConfiguration

	Start(s IService) (err error)
	Stop(s IService) (err error)

	// ReportStart method reports a successful start, ReportStop method reports
	// a successful stop. ReportStart method is called externally, i.e. by an
	// external entity. The ReportStop method, however, is called internally,
	// i.e. by the service component itself.
	ReportStart()
	ReportStop()
}
