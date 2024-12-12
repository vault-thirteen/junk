package s

import (
	"errors"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

// MessageSettings are settings of e-mail messages.
type MessageSettings struct {
	SubjectTemplateForRegVCode cmb.Text `json:"subjectTemplateForRegVCode"`
	SubjectTemplateForReg      cmb.Text `json:"subjectTemplateForReg"`
	BodyTemplateForRegVCode    cmb.Text `json:"bodyTemplateForRegVCode"`
	BodyTemplateForReg         cmb.Text `json:"bodyTemplateForReg"`
	BodyTemplateForLogIn       cmb.Text `json:"bodyTemplateForLogIn"`
	BodyTemplateForPwdChange   cmb.Text `json:"bodyTemplateForPwdChange"`
	BodyTemplateForEmailChange cmb.Text `json:"bodyTemplateForEmailChange"`
}

func (s MessageSettings) Check() (err error) {
	if (len(s.SubjectTemplateForRegVCode) == 0) ||
		(len(s.SubjectTemplateForReg) == 0) ||
		(len(s.BodyTemplateForRegVCode) == 0) ||
		(len(s.BodyTemplateForReg) == 0) ||
		(len(s.BodyTemplateForLogIn) == 0) ||
		(len(s.BodyTemplateForPwdChange) == 0) ||
		(len(s.BodyTemplateForEmailChange) == 0) {
		return errors.New(c.MsgMessageSettingError)
	}

	return nil
}
