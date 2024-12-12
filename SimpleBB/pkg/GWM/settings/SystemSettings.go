package s

import (
	"errors"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/ClientIPAddressSource"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

// systemSettings are system settings.
// Many of these settings must be synchronised with other modules.
type systemSettings struct {
	SettingsVersion base2.Count `json:"settingsVersion"`
	SiteName        base2.Text  `json:"siteName"`
	SiteDomain      base2.Text  `json:"siteDomain"`

	// Firewall.
	IsFirewallUsed base2.Flag `json:"isFirewallUsed"`

	// ClientIPAddressSource setting selects where to search for client's IP
	// address. '1' means that IP address is taken directly from the client's
	// address of the HTTP request; '2' means that IP address is taken from the
	// custom HTTP header which is configured by the ClientIPAddressHeader
	// setting. One of the most common examples of a custom header may be the
	// 'X-Forwarded-For' HTTP header. For most users the first variant ('1') is
	// the most suitable. The second variant ('2') may be used if you are
	// proxying requests of your clients somewhere inside your own network
	// infrastructure, such as via a load balancer or with a reverse proxy.
	ClientIPAddressSource derived1.IClientIPAddressSource `json:"clientIPAddressSource"`
	ClientIPAddressHeader string                          `json:"clientIPAddressHeader"`

	// Captcha.
	CaptchaImgServerHost string      `json:"captchaImgServerHost"`
	CaptchaImgServerPort uint16      `json:"captchaImgServerPort"`
	CaptchaFolder        simple.Path `json:"captchaFolder"`

	// Sessions and messages.
	SessionMaxDuration base2.Count `json:"sessionMaxDuration"`
	MessageEditTime    base2.Count `json:"messageEditTime"`
	PageSize           base2.Count `json:"pageSize"`

	// URL paths.
	ApiFolder              simple.Path `json:"apiFolder"`
	PublicSettingsFileName simple.Path `json:"publicSettingsFileName"`

	// Front end.
	IsFrontEndEnabled         base2.Flag  `json:"isFrontEndEnabled"`
	FrontEndStaticFilesFolder simple.Path `json:"frontEndStaticFilesFolder"`
	FrontEndAssetsFolder      simple.Path `json:"frontEndAssetsFolder"`

	// Development settings.
	IsDebugMode                               base2.Flag `json:"isDebugMode"`
	IsDeveloperMode                           base2.Flag `json:"isDeveloperMode"`
	DevModeHttpHeaderAccessControlAllowOrigin string     `json:"devModeHttpHeaderAccessControlAllowOrigin"`

	NotificationCountLimit base2.Count `json:"notificationCountLimit"`
}

func NewSystemSettings() (ss ISystemSettings) {
	return &systemSettings{
		ClientIPAddressSource: cm.NewClientIPAddressSource(),
	}
}

func (s systemSettings) Check() (err error) {
	if (s.SettingsVersion == 0) ||
		(len(s.SiteName) == 0) ||
		(len(s.SiteDomain) == 0) ||
		(s.ClientIPAddressSource.GetValue().RawValue() < cm.ClientIPAddressSource_Direct) ||
		(s.ClientIPAddressSource.GetValue().RawValue() > cm.ClientIPAddressSourceMax) ||
		(len(s.CaptchaImgServerHost) == 0) ||
		(s.CaptchaImgServerPort == 0) ||
		(len(s.CaptchaFolder) == 0) ||
		(s.SessionMaxDuration == 0) ||
		(s.MessageEditTime == 0) ||
		(s.PageSize == 0) ||
		(len(s.ApiFolder) == 0) ||
		(len(s.PublicSettingsFileName) == 0) ||
		(s.NotificationCountLimit == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	if s.IsFrontEndEnabled {
		if (len(s.FrontEndStaticFilesFolder) == 0) ||
			(len(s.FrontEndAssetsFolder) == 0) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	if s.ClientIPAddressSource.GetValue().RawValue() == cm.ClientIPAddressSource_CustomHeader {
		if len(s.ClientIPAddressHeader) == 0 {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	if s.IsDeveloperMode {
		if len(s.DevModeHttpHeaderAccessControlAllowOrigin) == 0 {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}

// Emulated class members.
func (s systemSettings) GetSettingsVersion() base2.Count { return s.SettingsVersion }
func (s systemSettings) GetSiteName() base2.Text         { return s.SiteName }
func (s systemSettings) GetSiteDomain() base2.Text       { return s.SiteDomain }
func (s systemSettings) GetIsFirewallUsed() base2.Flag   { return s.IsFirewallUsed }
func (s systemSettings) GetClientIPAddressSource() derived1.IClientIPAddressSource {
	return s.ClientIPAddressSource
}
func (s systemSettings) GetClientIPAddressHeader() string       { return s.ClientIPAddressHeader }
func (s systemSettings) GetCaptchaImgServerHost() string        { return s.CaptchaImgServerHost }
func (s systemSettings) GetCaptchaImgServerPort() uint16        { return s.CaptchaImgServerPort }
func (s systemSettings) GetCaptchaFolder() simple.Path          { return s.CaptchaFolder }
func (s systemSettings) GetSessionMaxDuration() base2.Count     { return s.SessionMaxDuration }
func (s systemSettings) GetMessageEditTime() base2.Count        { return s.MessageEditTime }
func (s systemSettings) GetPageSize() base2.Count               { return s.PageSize }
func (s systemSettings) GetApiFolder() simple.Path              { return s.ApiFolder }
func (s systemSettings) GetPublicSettingsFileName() simple.Path { return s.PublicSettingsFileName }
func (s systemSettings) GetIsFrontEndEnabled() base2.Flag       { return s.IsFrontEndEnabled }
func (s systemSettings) GetFrontEndStaticFilesFolder() simple.Path {
	return s.FrontEndStaticFilesFolder
}
func (s systemSettings) GetFrontEndAssetsFolder() simple.Path { return s.FrontEndAssetsFolder }
func (s systemSettings) GetIsDebugMode() base2.Flag           { return s.IsDebugMode }
func (s systemSettings) GetIsDeveloperMode() base2.Flag       { return s.IsDeveloperMode }
func (s systemSettings) GetDevModeHttpHeaderAccessControlAllowOrigin() string {
	return s.DevModeHttpHeaderAccessControlAllowOrigin
}
func (s systemSettings) GetNotificationCountLimit() base2.Count { return s.NotificationCountLimit }
