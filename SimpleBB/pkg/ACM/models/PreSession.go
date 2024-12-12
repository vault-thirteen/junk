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

type PreSession struct {
	Id             base2.Id
	UserId         base2.Id
	TimeOfCreation time.Time
	RequestId      simple.RequestId

	// IP address of a user. B = Byte array.
	UserIPAB net.IP

	AuthDataBytes        cmr.AuthChallengeData
	IsCaptchaRequired    base2.Flag
	CaptchaId            *simple.CaptchaId
	IsVerifiedByCaptcha  *base2.Flag
	IsVerifiedByPassword base2.Flag

	// Verification code is set on Step 2, so it is NULL on Step 1.
	VerificationCode *simple.VerificationCode

	IsEmailSent       base2.Flag
	IsVerifiedByEmail base2.Flag
}

func NewPreSession() (ps *PreSession) {
	return &PreSession{}
}

func NewPreSessionFromScannableSource(src cmi.IScannable) (ps *PreSession, err error) {
	ps = NewPreSession()

	err = src.Scan(
		&ps.Id,
		&ps.UserId,
		&ps.TimeOfCreation,
		&ps.RequestId,
		&ps.UserIPAB,
		&ps.AuthDataBytes,
		&ps.IsCaptchaRequired,
		&ps.CaptchaId,
		&ps.IsVerifiedByCaptcha,
		&ps.IsVerifiedByPassword,
		&ps.VerificationCode,
		&ps.IsEmailSent,
		&ps.IsVerifiedByEmail,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ps, nil
}
