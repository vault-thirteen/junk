package c

import (
	"sync"

	"github.com/vault-thirteen/TR1/src/components/MailerComponent"
	"github.com/vault-thirteen/TR1/src/libraries/mailer"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationServiceEntry"
)

type ControllerFastAccessRegistry struct {
	systemSettings *ccse.CommonConfigurationServiceEntry
	mc             *mc.MailerComponent
	m              *mailer.Mailer
	mg             *sync.Mutex
}
