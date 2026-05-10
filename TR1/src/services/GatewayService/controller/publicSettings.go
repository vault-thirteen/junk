package c

import (
	"encoding/json"
	"net/http"

	cm "github.com/vault-thirteen/TR1/src/models/common"
	hm "github.com/vault-thirteen/TR1/src/models/http"
	ccp "github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
	"github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

func (c *Controller) initPublicSettings() (err error) {
	ps := &cm.PublicSettings{
		Version:            c.far.systemSettings.GetParameterAsString(ccp.PublicSettingsVersion),
		TTL:                c.far.systemSettings.GetParameterAsInt(ccp.PublicSettingsTtl),
		SiteName:           c.far.systemSettings.GetParameterAsString(ccp.SiteName),
		SiteDomain:         c.far.systemSettings.GetParameterAsString(ccp.SiteDomain),
		SessionMaxDuration: c.far.sessionMaxDuration,
		MessageEditTime:    c.far.messageEditTime,
		PageSize:           c.far.pageSize,
	}

	c.publicSettingsFile, err = json.Marshal(ps)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) handleSettingsRequest(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		c.respondMethodNotAllowed(rw)
		return
	}

	// Check accepted MIME types.
	ok, err := hh.CheckBrowserSupportForJson(req)
	if err != nil {
		c.respondBadRequest(rw)
		return
	}
	if !ok {
		c.respondNotAcceptable(rw)
		return
	}

	rw.Header().Set(header.HttpHeaderContentType, hm.ContentType_Json)

	_, err = rw.Write(c.publicSettingsFile)
	if err != nil {
		c.logError(err)
		return
	}
}
