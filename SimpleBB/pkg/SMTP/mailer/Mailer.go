package mailer

import (
	"errors"

	"github.com/valord577/mailx"
)

// This tiny library uses another SMTP library because the built-in Golang SMTP
// library is so poor that it hangs infinitely during the connection to an SMTP
// server.

const (
	ErrHostIsEmpty      = "host is empty"
	ErrPortIsNotSet     = "port is not set"
	ErrUsernameIsEmpty  = "username is empty"
	ErrUserAgentIsEmpty = "user agent is empty"
)

// Mailer is an e-mail sender.
// It sends e-mail messages using the SMTP protocol.
type Mailer struct {
	host      string
	port      uint16
	username  string
	pwd       string
	userAgent string
}

func NewMailer(
	host string,
	port uint16,
	username string,
	pwd string,
	userAgent string,
) (m *Mailer, err error) {
	if len(host) == 0 {
		return nil, errors.New(ErrHostIsEmpty)
	}
	if port == 0 {
		return nil, errors.New(ErrPortIsNotSet)
	}
	if len(username) == 0 {
		return nil, errors.New(ErrUsernameIsEmpty)
	}
	if len(userAgent) == 0 {
		return nil, errors.New(ErrUserAgentIsEmpty)
	}

	return &Mailer{
		host:      host,
		port:      port,
		username:  username,
		pwd:       pwd,
		userAgent: userAgent,
	}, nil
}

func (m *Mailer) SendMail(
	recipients []string,
	subject string,
	message string,
) (err error) {
	msg := mailx.NewMessage()
	msg.SetTo(recipients...)
	msg.SetSubject(subject)
	msg.SetPlainBody(message)
	msg.SetUserAgent(m.userAgent)

	dialer := mailx.Dialer{
		Host:         m.host,
		Port:         int(m.port),
		Username:     m.username,
		Password:     m.pwd,
		SSLOnConnect: true,
	}

	err = dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
