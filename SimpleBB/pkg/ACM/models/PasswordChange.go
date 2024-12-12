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

type PasswordChange struct {
	Id             base2.Id
	UserId         base2.Id
	TimeOfCreation time.Time
	RequestId      *simple.RequestId

	// IP address of a user. B = Byte array.
	UserIPAB net.IP

	AuthDataBytes        cmr.AuthChallengeData
	IsCaptchaRequired    base2.Flag
	CaptchaId            *simple.CaptchaId
	IsVerifiedByCaptcha  *base2.Flag
	IsVerifiedByPassword base2.Flag
	VerificationCode     *simple.VerificationCode
	IsEmailSent          base2.Flag
	IsVerifiedByEmail    base2.Flag
	NewPasswordBytes     []byte
}

type PasswordChangeVerificationFlags struct {
	IsVerifiedByCaptcha  *base2.Flag
	IsVerifiedByPassword base2.Flag
	IsVerifiedByEmail    base2.Flag
}

func NewPasswordChange() (pc *PasswordChange) {
	return &PasswordChange{}
}

func NewPasswordChangeFromScannableSource(src cmi.IScannable) (pc *PasswordChange, err error) {
	pc = NewPasswordChange()

	err = src.Scan(
		&pc.Id,
		&pc.UserId,
		&pc.TimeOfCreation,
		&pc.RequestId,
		&pc.UserIPAB,
		&pc.AuthDataBytes,
		&pc.IsCaptchaRequired,
		&pc.CaptchaId,
		&pc.IsVerifiedByCaptcha,
		&pc.IsVerifiedByPassword,
		&pc.VerificationCode,
		&pc.IsEmailSent,
		&pc.IsVerifiedByEmail,
		&pc.NewPasswordBytes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return pc, nil
}
