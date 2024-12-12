package km

import (
	"crypto/rsa"
	"errors"
	"fmt"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	WebTokenField_UserId    = "userId"
	WebTokenField_SessionId = "sessionId"
	TokenHeader_Alg         = "alg"
	TokenAlg_PS512          = "PS512" // RSA-PSS.
	TokenAlg_RS512          = "RS512" // RSA.
)

const (
	ErrSigningMethodIsNotSupported = "signing method is not supported"
	ErrTokenIsNotValid             = "token is not valid"
	ErrTokenIsBroken               = "token is broken"
	ErrFUnsupportedSigningMethod   = "unsupported signing method: %s"
	ErrFUnexpectedSigningMethod    = "unexpected signing method: %v"
	ErrTypeCast                    = "type cast error"
)

type KeyMaker struct {
	privateKey        *rsa.PrivateKey
	publicKey         *rsa.PublicKey
	signingMethod     jwt.SigningMethod
	signingMethodName string
}

func New(signingMethodName string, privateKeyFilePath simple.Path, publicKeyFilePath simple.Path) (km *KeyMaker, err error) {
	var signingMethod jwt.SigningMethod
	switch signingMethodName {
	case TokenAlg_PS512:
		signingMethod = jwt.SigningMethodPS512
	case TokenAlg_RS512:
		signingMethod = jwt.SigningMethodRS512
	default:
		return nil, errors.New(ErrSigningMethodIsNotSupported)
	}

	km = &KeyMaker{
		signingMethod:     signingMethod,
		signingMethodName: strings.ToUpper(signingMethodName),
	}

	km.privateKey, err = getPrivateKey(privateKeyFilePath, signingMethodName)
	if err != nil {
		return nil, err
	}

	km.publicKey, err = getPublicKey(publicKeyFilePath, signingMethodName)
	if err != nil {
		return nil, err
	}

	return km, nil
}

func (km *KeyMaker) MakeJWToken(userId cmb.Id, sessionId cmb.Id) (tokenString simple.WebTokenString, err error) {
	claims := jwt.MapClaims{
		WebTokenField_UserId:    userId,
		WebTokenField_SessionId: sessionId,
	}

	token := jwt.NewWithClaims(km.signingMethod, claims, nil)

	var s string
	s, err = token.SignedString(km.privateKey)
	if err != nil {
		return "", err
	}

	return simple.WebTokenString(s), nil
}

func (km *KeyMaker) ValidateToken(tokenString simple.WebTokenString) (userId cmb.Id, sessionId cmb.Id, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString.ToString(), func(token *jwt.Token) (interface{}, error) {
		if strings.ToUpper(token.Method.Alg()) != km.signingMethodName {
			return nil, fmt.Errorf(ErrFUnsupportedSigningMethod, token.Method.Alg())
		}

		tmp := token.Header[TokenHeader_Alg]
		algString := tmp.(string)
		if strings.ToUpper(algString) != km.signingMethodName {
			return nil, fmt.Errorf(ErrFUnsupportedSigningMethod, token.Header[TokenHeader_Alg])
		}

		var ok bool
		switch km.signingMethodName {
		case TokenAlg_PS512:
			_, ok = token.Method.(*jwt.SigningMethodRSAPSS)
		case TokenAlg_RS512:
			_, ok = token.Method.(*jwt.SigningMethodRSA)
		}
		if !ok {
			return nil, fmt.Errorf(ErrFUnexpectedSigningMethod, token.Header[TokenHeader_Alg])
		}

		return km.publicKey, nil
	})
	if err != nil {
		return 0, 0, err
	}

	if !token.Valid {
		return 0, 0, errors.New(ErrTokenIsNotValid)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, errors.New(ErrTokenIsBroken)
	}

	var userIdIfc any
	userIdIfc, ok = claims[WebTokenField_UserId]
	if !ok {
		return 0, 0, errors.New(ErrTokenIsBroken)
	}

	var userIdFloat64 float64
	userIdFloat64, ok = userIdIfc.(float64)
	if !ok {
		return 0, 0, errors.New(ErrTokenIsBroken)
	}

	var sessionIdIfc any
	sessionIdIfc, ok = claims[WebTokenField_SessionId]
	if !ok {
		return 0, 0, errors.New(ErrTokenIsBroken)
	}

	var sessionIdFloat64 float64
	sessionIdFloat64, ok = sessionIdIfc.(float64)
	if !ok {
		return 0, 0, errors.New(ErrTokenIsBroken)
	}

	return cmb.Id(userIdFloat64), cmb.Id(sessionIdFloat64), nil
}
