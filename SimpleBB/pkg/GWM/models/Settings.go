package models

import (
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

type Settings struct {
	Version                   base2.Count `json:"version"`
	ProductVersion            base2.Text  `json:"productVersion"`
	SiteName                  base2.Text  `json:"siteName"`
	SiteDomain                base2.Text  `json:"siteDomain"`
	CaptchaFolder             cm.Path     `json:"captchaFolder"`
	SessionMaxDuration        base2.Count `json:"sessionMaxDuration"`
	MessageEditTime           base2.Count `json:"messageEditTime"`
	PageSize                  base2.Count `json:"pageSize"`
	ApiFolder                 cm.Path     `json:"apiFolder"`
	PublicSettingsFileName    cm.Path     `json:"publicSettingsFileName"`
	IsFrontEndEnabled         base2.Flag  `json:"isFrontEndEnabled"`
	FrontEndStaticFilesFolder cm.Path     `json:"frontEndStaticFilesFolder"`
	NotificationCountLimit    base2.Count `json:"notificationCountLimit"`
}
