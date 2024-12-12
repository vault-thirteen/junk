package app

import (
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Module"
)

const (
	ServiceName_ACM  = "Access Control Module"
	ServiceName_GWM  = "Gateway Module"
	ServiceName_MM   = "Message Module"
	ServiceName_NM   = "Notification Module"
	ServiceName_RCS  = "Captcha Module"
	ServiceName_SM   = "Subscription Module"
	ServiceName_SMTP = "SMTP Module"
)

const (
	ConfigurationFilePathDefault_ACM  = "ACM.json"
	ConfigurationFilePathDefault_GWM  = "GWM.json"
	ConfigurationFilePathDefault_MM   = "MM.json"
	ConfigurationFilePathDefault_NM   = "NM.json"
	ConfigurationFilePathDefault_RCS  = "RCS.json"
	ConfigurationFilePathDefault_SM   = "SM.json"
	ConfigurationFilePathDefault_SMTP = "SMTP.json"
)

const (
	ServiceShortName_ACM  = "ACM"
	ServiceShortName_GWM  = "GWM"
	ServiceShortName_MM   = "MM"
	ServiceShortName_NM   = "NM"
	ServiceShortName_RCS  = "RCS"
	ServiceShortName_SM   = "SM"
	ServiceShortName_SMTP = "SMTP"
)

const (
	ModuleId_ACM  = cm.Module_ACM
	ModuleId_GWM  = cm.Module_GWM
	ModuleId_MM   = cm.Module_MM
	ModuleId_NM   = cm.Module_NM
	ModuleId_RCS  = cm.Module_RCS
	ModuleId_SM   = cm.Module_SM
	ModuleId_SMTP = cm.Module_SMTP
)
