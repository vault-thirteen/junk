package mc

import (
	"fmt"
	"sync"

	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/libraries/mailer"
	"github.com/vault-thirteen/TR1/src/models/common"
	"github.com/vault-thirteen/TR1/src/shared/CommonConfigurationParameter"
)

type MailerComponent struct {
	cfg         interfaces.IConfiguration
	mailer      *mailer.Mailer
	mailerGuard sync.Mutex
}

func (c *MailerComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	mc := &MailerComponent{
		cfg: cfg,
	}

	mailerSettings := cfg.GetComponent(cm.Component_Mailer, cm.Protocol_None)
	host := mailerSettings.GetParameterAsString(ccp.Host)
	port := mailerSettings.GetParameterAsInt(ccp.Port)
	user := mailerSettings.GetParameterAsString(ccp.User)

	password := mailerSettings.GetParameterAsString(ccp.Password)
	if len(password) == 0 {
		password, err = cm.GetPasswordFromStdin("SMTP server")
		if err != nil {
			return nil, err
		}
	}

	userAgent := mailerSettings.GetParameterAsString(ccp.UserAgent)

	mc.mailer, err = mailer.NewMailer(host, uint16(port), user, password, userAgent)
	if err != nil {
		return nil, err
	}

	return mc, nil
}
func (c *MailerComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *MailerComponent) Start(s interfaces.IService) (err error) {
	return nil
}
func (c *MailerComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	c.ReportStop()

	return nil
}

func (c *MailerComponent) ReportStart() {
	fmt.Println(fmt.Sprintf("MailerComponent has started"))
}
func (c *MailerComponent) ReportStop() {
	fmt.Println("MailerComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *MailerComponent) {
	return x.(*MailerComponent)
}

// Non-standard methods.

func (c *MailerComponent) GetMailer() (mailer *mailer.Mailer) {
	return c.mailer
}
func (c *MailerComponent) GetMailerGuard() (mailerGuard *sync.Mutex) {
	return &c.mailerGuard
}
