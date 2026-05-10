package km

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	signingMethodName string
	publicKey         *rsa.PublicKey

	claims    jwt.MapClaims
	userId    int
	sessionId int
	expTime   int64
}

func NewValidator(signingMethodName string, publicKey *rsa.PublicKey) *Validator {
	return &Validator{
		signingMethodName: signingMethodName,
		publicKey:         publicKey,
	}
}

func (v *Validator) KeyFn(token *jwt.Token) (key interface{}, err error) {
	if strings.ToUpper(token.Method.Alg()) != v.signingMethodName {
		return nil, fmt.Errorf(ErrFUnsupportedSigningMethod, token.Method.Alg())
	}

	tmp := token.Header[TokenHeader_Alg]
	algString := tmp.(string)
	if strings.ToUpper(algString) != v.signingMethodName {
		return nil, fmt.Errorf(ErrFUnsupportedSigningMethod, token.Header[TokenHeader_Alg])
	}

	var ok bool
	switch v.signingMethodName {
	case TokenAlg_PS512:
		_, ok = token.Method.(*jwt.SigningMethodRSAPSS)
	case TokenAlg_RS512:
		_, ok = token.Method.(*jwt.SigningMethodRSA)
	}
	if !ok {
		return nil, fmt.Errorf(ErrFUnexpectedSigningMethod, token.Header[TokenHeader_Alg])
	}

	v.claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New(ErrTokenIsBroken)
	}

	err = v.getUserId()
	if err != nil {
		return nil, err
	}

	err = v.getSessionId()
	if err != nil {
		return nil, err
	}

	err = v.getExpTime()
	if err != nil {
		return nil, err
	}

	return v.publicKey, nil
}

func (v *Validator) getUserId() (err error) {
	v.userId, err = v.getClaimAsInt(WebTokenField_UserId)
	if err != nil {
		return err
	}

	return nil
}
func (v *Validator) getSessionId() (err error) {
	v.sessionId, err = v.getClaimAsInt(WebTokenField_SessionId)
	if err != nil {
		return err
	}

	return nil
}
func (v *Validator) getExpTime() (err error) {
	var x int
	x, err = v.getClaimAsInt(WebTokenField_ExpirationTime)
	if err != nil {
		return err
	}

	v.expTime = int64(x)

	return nil
}
func (v *Validator) getClaimAsInt(claimName string) (claim int, err error) {
	c, ok := v.claims[claimName]
	if !ok {
		return 0, errors.New(ErrTokenIsBroken)
	}

	var cAsFloat64 float64
	cAsFloat64, ok = c.(float64)
	if !ok {
		return 0, errors.New(ErrTokenIsBroken)
	}

	return int(cAsFloat64), nil
}
