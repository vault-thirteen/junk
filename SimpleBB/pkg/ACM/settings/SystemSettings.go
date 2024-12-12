package s

import (
	"errors"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

// SystemSettings are system settings.
type SystemSettings struct {
	SiteName                     base2.Text  `json:"siteName"`
	SiteDomain                   base2.Text  `json:"siteDomain"`
	VerificationCodeLength       base2.Count `json:"verificationCodeLength"`
	UserNameMaxLenInBytes        base2.Count `json:"userNameMaxLenInBytes"`
	UserPasswordMaxLenInBytes    base2.Count `json:"userPasswordMaxLenInBytes"`
	PreRegUserExpirationTime     base2.Count `json:"preRegUserExpirationTime"`
	IsAdminApprovalRequired      base2.Flag  `json:"isAdminApprovalRequired"`
	LogInRequestIdLength         base2.Count `json:"logInRequestIdLength"`
	LogInTryTimeout              base2.Count `json:"logInTryTimeout"`
	PreSessionExpirationTime     base2.Count `json:"preSessionExpirationTime"`
	SessionMaxDuration           base2.Count `json:"sessionMaxDuration"`
	PasswordChangeExpirationTime base2.Count `json:"passwordChangeExpirationTime"`
	EmailChangeExpirationTime    base2.Count `json:"emailChangeExpirationTime"`
	ActionTryTimeout             base2.Count `json:"actionTryTimeout"`
	PageSize                     base2.Count `json:"pageSize"`

	// This setting must be synchronised with settings of the Gateway module.
	IsTableOfIncidentsUsed base2.Flag `json:"isTableOfIncidentsUsed"`

	// This setting is used only when a table of incidents is enabled.
	BlockTimePerIncident BlockTimePerIncident `json:"blockTimePerIncident"`

	IsDebugMode bool `json:"isDebugMode"`
}

// BlockTimePerIncident is block time in seconds for each type of incident.
type BlockTimePerIncident struct {
	IllegalAccessAttempt     base2.Count `json:"illegalAccessAttempt"`     // 1.
	FakeToken                base2.Count `json:"fakeToken"`                // 2.
	VerificationCodeMismatch base2.Count `json:"verificationCodeMismatch"` // 3.
	DoubleLogInAttempt       base2.Count `json:"doubleLogInAttempt"`       // 4.
	PreSessionHacking        base2.Count `json:"preSessionHacking"`        // 5.
	CaptchaAnswerMismatch    base2.Count `json:"captchaAnswerMismatch"`    // 6.
	PasswordMismatch         base2.Count `json:"passwordMismatch"`         // 7.
	PasswordChangeHacking    base2.Count `json:"passwordChangeHacking"`    // 8.
	EmailChangeHacking       base2.Count `json:"emailChangeHacking"`       // 9.
	FakeIPA                  base2.Count `json:"fakeIPA"`                  // 10.
}

func (s SystemSettings) Check() (err error) {
	if (len(s.SiteName) == 0) ||
		(len(s.SiteDomain) == 0) ||
		(s.VerificationCodeLength < 8) ||
		(s.UserNameMaxLenInBytes == 0) ||
		(s.UserPasswordMaxLenInBytes == 0) ||
		(s.PreRegUserExpirationTime == 0) ||
		(s.LogInRequestIdLength == 0) ||
		(s.LogInTryTimeout == 0) ||
		(s.PreSessionExpirationTime == 0) ||
		(s.SessionMaxDuration == 0) ||
		(s.PasswordChangeExpirationTime == 0) ||
		(s.EmailChangeExpirationTime == 0) ||
		(s.ActionTryTimeout == 0) ||
		(s.PageSize == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	// Firewall.
	if s.IsTableOfIncidentsUsed {
		if (s.BlockTimePerIncident.IllegalAccessAttempt == 0) ||
			(s.BlockTimePerIncident.FakeToken == 0) ||
			(s.BlockTimePerIncident.VerificationCodeMismatch == 0) ||
			(s.BlockTimePerIncident.DoubleLogInAttempt == 0) ||
			(s.BlockTimePerIncident.PreSessionHacking == 0) ||
			(s.BlockTimePerIncident.CaptchaAnswerMismatch == 0) ||
			(s.BlockTimePerIncident.PasswordMismatch == 0) ||
			(s.BlockTimePerIncident.PasswordChangeHacking == 0) ||
			(s.BlockTimePerIncident.EmailChangeHacking == 0) ||
			(s.BlockTimePerIncident.FakeIPA == 0) {
			return errors.New(c.MsgSystemSettingError)
		}
	}

	return nil
}
