package c

import (
	rcs "github.com/vault-thirteen/RingCaptcha/server"
	"github.com/vault-thirteen/TR1/src/components/CaptchaComponent"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
)

type ControllerFastAccessRegistry struct {
	systemSettings *ccse.CommonConfigurationServiceEntry
	cc             *cc.CaptchaComponent
	cs             *rcs.CaptchaServer
}
