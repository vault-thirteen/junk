package s

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

type ISystemSettings interface {
	Check() (err error)

	// Emulated class members.
	GetSettingsVersion() base2.Count
	GetSiteName() base2.Text
	GetSiteDomain() base2.Text
	GetIsFirewallUsed() base2.Flag
	GetClientIPAddressSource() derived1.IClientIPAddressSource
	GetClientIPAddressHeader() string
	GetCaptchaImgServerHost() string
	GetCaptchaImgServerPort() uint16
	GetCaptchaFolder() simple.Path
	GetSessionMaxDuration() base2.Count
	GetMessageEditTime() base2.Count
	GetPageSize() base2.Count
	GetApiFolder() simple.Path
	GetPublicSettingsFileName() simple.Path
	GetIsFrontEndEnabled() base2.Flag
	GetFrontEndStaticFilesFolder() simple.Path
	GetFrontEndAssetsFolder() simple.Path
	GetIsDebugMode() base2.Flag
	GetIsDeveloperMode() base2.Flag
	GetDevModeHttpHeaderAccessControlAllowOrigin() string
	GetNotificationCountLimit() base2.Count
}
