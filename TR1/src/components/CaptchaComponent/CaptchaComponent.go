package cc

import (
	"fmt"

	rcm "github.com/vault-thirteen/RingCaptcha/models"
	rcs "github.com/vault-thirteen/RingCaptcha/server"
	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

const (
	Err_UnexpectedResponse = "unexpected response"
)

type CaptchaComponent struct {
	cfg           interfaces.IConfiguration
	errorsChan    *chan error
	captchaServer *rcs.CaptchaServer
}

func (c *CaptchaComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	cc := &CaptchaComponent{
		cfg:        cfg,
		errorsChan: controller.GetErrorsChan(),
	}

	captchaSettings := cfg.GetComponent(cm.Component_Captcha, cm.Protocol_None)

	imageServerSettings := cfg.GetServer(cm.ServerType_External, cm.Protocol_HTTP)
	imageServerHost := imageServerSettings.GetParameterAsString(ccp.Host)
	imageServerPort := imageServerSettings.GetParameterAsInt(ccp.Port)
	imageServerName := imageServerSettings.GetParameterAsString(ccp.Name)

	var css = &rcm.CaptchaServerSettings{
		// Main settings.
		IsImageStorageUsed:        captchaSettings.GetParameterAsBool(ccp.IsImageStorageUsed),
		IsImageServerEnabled:      captchaSettings.GetParameterAsBool(ccp.IsImageServerEnabled),
		IsImageCleanupAtStartUsed: captchaSettings.GetParameterAsBool(ccp.IsImageCleanupAtStartUsed),
		IsStorageCleaningEnabled:  captchaSettings.GetParameterAsBool(ccp.IsStorageCleaningEnabled),

		// Image settings.
		ImagesFolder:      captchaSettings.GetParameterAsString(ccp.ImagesFolder),
		ImageWidth:        uint(captchaSettings.GetParameterAsInt(ccp.ImageWidth)),
		ImageHeight:       uint(captchaSettings.GetParameterAsInt(ccp.ImageHeight)),
		FilesCountToClean: captchaSettings.GetParameterAsInt(ccp.FilesCountToClean),

		// HTTP server settings.
		HttpHost:       imageServerHost,
		HttpPort:       uint16(imageServerPort),
		HttpErrorsChan: cc.errorsChan,
		HttpServerName: imageServerName,

		// File cache settings.
		FileCacheSizeLimit:   captchaSettings.GetParameterAsInt(ccp.FileCacheSizeLimit),
		FileCacheVolumeLimit: captchaSettings.GetParameterAsInt(ccp.FileCacheVolumeLimit),
		FileCacheItemTtl:     uint(captchaSettings.GetParameterAsInt(ccp.FileCacheItemTtl)),

		// Record cache settings.
		RecordCacheSizeLimit: captchaSettings.GetParameterAsInt(ccp.RecordCacheSizeLimit),
		RecordCacheItemTtl:   uint(captchaSettings.GetParameterAsInt(ccp.RecordCacheItemTtl)),
	}

	cc.captchaServer, err = rcs.NewCaptchaServer(css)
	if err != nil {
		return nil, err
	}

	return cc, nil
}
func (c *CaptchaComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *CaptchaComponent) Start(s interfaces.IService) (err error) {
	err = c.captchaServer.Start()
	if err != nil {
		return err
	}

	return nil
}
func (c *CaptchaComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	err = c.captchaServer.Stop()
	if err != nil {
		return err
	}
	c.ReportStop()

	return nil
}

func (c *CaptchaComponent) ReportStart() {
	fmt.Println(fmt.Sprintf("CaptchaComponent has started"))
}
func (c *CaptchaComponent) ReportStop() {
	fmt.Println("CaptchaComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *CaptchaComponent) {
	return x.(*CaptchaComponent)
}

// Non-standard methods.

func (c *CaptchaComponent) GetCaptchaServer() (captchaServer *rcs.CaptchaServer) {
	return c.captchaServer
}
