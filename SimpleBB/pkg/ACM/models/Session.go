package models

import (
	"database/sql"
	"errors"
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"net"
	"time"
)

type Session struct {
	Id        cmb.Id    `json:"id"`
	UserId    cmb.Id    `json:"userId"`
	StartTime time.Time `json:"startTime"`

	// IP address of a user. B = Byte array.
	UserIPAB net.IP `json:"userIPA"`
}

func NewSession() (s *Session) {
	return &Session{}
}

func NewSessionFromScannableSource(src cmi.IScannable) (s *Session, err error) {
	s = NewSession()

	err = src.Scan(
		&s.Id,
		&s.UserId,
		&s.StartTime,
		&s.UserIPAB,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return s, nil
}
