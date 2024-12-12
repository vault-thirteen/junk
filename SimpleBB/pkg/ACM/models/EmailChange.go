package models

import (
	"database/sql"
	"errors"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"net"
	"time"
)

type EmailChange struct {
	Id             base2.Id
	UserId         base2.Id
	TimeOfCreation time.Time
	RequestId      *simple.RequestId

	// IP address of a user. B = Byte array.
	UserIPAB net.IP

	AuthDataBytes     cmr.AuthChallengeData
	IsCaptchaRequired base2.Flag
	CaptchaId         *simple.CaptchaId
	EmailChangeVerificationFlags

	// Old e-mail.
	VerificationCodeOld *simple.VerificationCode
	IsOldEmailSent      base2.Flag

	// New e-mail.
	NewEmail            simple.Email
	VerificationCodeNew *simple.VerificationCode
	IsNewEmailSent      base2.Flag
}

type EmailChangeVerificationFlags struct {
	IsVerifiedByCaptcha  *base2.Flag
	IsVerifiedByPassword base2.Flag
	IsVerifiedByOldEmail base2.Flag
	IsVerifiedByNewEmail base2.Flag
}

func NewEmailChange() (ec *EmailChange) {
	return &EmailChange{}
}

func NewEmailChangeFromScannableSource(src cmi.IScannable) (ec *EmailChange, err error) {
	ec = NewEmailChange()

	err = src.Scan(
		&ec.Id,
		&ec.UserId,
		&ec.TimeOfCreation,
		&ec.RequestId,
		&ec.UserIPAB,
		&ec.AuthDataBytes,
		&ec.IsCaptchaRequired,
		&ec.CaptchaId,
		&ec.IsVerifiedByCaptcha,
		&ec.IsVerifiedByPassword,
		&ec.VerificationCodeOld,
		&ec.IsOldEmailSent,
		&ec.IsVerifiedByOldEmail,
		&ec.NewEmail,
		&ec.VerificationCodeNew,
		&ec.IsNewEmailSent,
		&ec.IsVerifiedByNewEmail,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ec, nil
}
