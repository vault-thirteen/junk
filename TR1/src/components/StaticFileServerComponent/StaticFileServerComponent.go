package sfsc

import (
	"fmt"

	"github.com/vault-thirteen/Simple-File-Server"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

type StaticFileServerComponent struct {
	cfg        interfaces.IConfiguration
	fileServer *sfs.SimpleFileServer
}

func (c *StaticFileServerComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	sfsc := &StaticFileServerComponent{
		cfg: cfg,
	}

	sfsSettings := cfg.GetComponent(cm.Component_SFS, cm.Protocol_None)

	sfsc.fileServer, err = sfs.NewSimpleFileServer(
		sfsSettings.GetParameterAsString(ccp.RootFolderPath),
		[]string{},
		sfsSettings.GetParameterAsBool(ccp.IsCacheEnabled),
		sfsSettings.GetParameterAsInt(ccp.FileCacheSizeLimit),
		sfsSettings.GetParameterAsInt(ccp.FileCacheVolumeLimit),
		uint(sfsSettings.GetParameterAsInt(ccp.CacheRecordTtl)),
	)
	if err != nil {
		return nil, err
	}

	return sfsc, nil
}
func (c *StaticFileServerComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *StaticFileServerComponent) Start(s interfaces.IService) (err error) {
	return nil
}
func (c *StaticFileServerComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	c.ReportStop()

	return nil
}

func (c *StaticFileServerComponent) ReportStart() {
	fmt.Println(fmt.Sprintf("StaticFileServerComponent has started"))
}
func (c *StaticFileServerComponent) ReportStop() {
	fmt.Println("StaticFileServerComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *StaticFileServerComponent) {
	return x.(*StaticFileServerComponent)
}

// Non-standard methods.

func (c *StaticFileServerComponent) GetFileServer() (fileServer *sfs.SimpleFileServer) {
	return c.fileServer
}
