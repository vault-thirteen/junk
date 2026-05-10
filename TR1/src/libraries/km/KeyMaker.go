package km

import (
	"crypto/rsa"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	nvl "github.com/vault-thirteen/Cache/NVL"
)

const (
	WebTokenField_UserId         = "userId"
	WebTokenField_SessionId      = "sessionId"
	WebTokenField_ExpirationTime = "exp"
	TokenHeader_Alg              = "alg"
	TokenAlg_PS512               = "PS512" // RSA-PSS.
	TokenAlg_RS512               = "RS512" // RSA.
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
	settings          *KeyMakerSettings
	privateKey        *rsa.PrivateKey
	publicKey         *rsa.PublicKey
	signingMethod     jwt.SigningMethod
	signingMethodName string
	cache             *nvl.Cache[string, *CachedData]
}

func New(settings *KeyMakerSettings) (km *KeyMaker, err error) {
	var signingMethod jwt.SigningMethod
	switch settings.SigningMethodName {
	case TokenAlg_PS512:
		signingMethod = jwt.SigningMethodPS512
	case TokenAlg_RS512:
		signingMethod = jwt.SigningMethodRS512
	default:
		return nil, errors.New(ErrSigningMethodIsNotSupported)
	}

	km = &KeyMaker{
		settings:          settings,
		signingMethod:     signingMethod,
		signingMethodName: strings.ToUpper(settings.SigningMethodName),
	}

	if settings.IsCacheEnabled {
		km.cache = nvl.NewCache[string, *CachedData](settings.CacheSizeLimit, settings.CacheRecordTtl)
	}

	km.privateKey, err = getPrivateKey(settings.PrivateKeyFilePath, settings.SigningMethodName)
	if err != nil {
		return nil, err
	}

	km.publicKey, err = getPublicKey(settings.PublicKeyFilePath, settings.SigningMethodName)
	if err != nil {
		return nil, err
	}

	return km, nil
}

func (km *KeyMaker) MakeJWToken(userId int, sessionId int, expirationTime time.Time) (tokenString string, err error) {
	claims := jwt.MapClaims{
		WebTokenField_UserId:         userId,
		WebTokenField_SessionId:      sessionId,
		WebTokenField_ExpirationTime: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(km.signingMethod, claims, nil)

	var s string
	s, err = token.SignedString(km.privateKey)
	if err != nil {
		return "", err
	}

	return s, nil
}

func (km *KeyMaker) ValidateToken(tokenString string) (userId int, sessionId int, err error) {
	// No cache is used.
	if !km.settings.IsCacheEnabled {
		var validator *Validator
		validator, err = km.validateToken(tokenString)
		if err != nil {
			return validator.userId, validator.sessionId, err
		}
		return validator.userId, validator.sessionId, nil
	}

	// Cache is used and data is ready.
	var cd *CachedData
	cd = km.getCachedData(tokenString)
	if cd != nil {
		return cd.UserId, cd.SessionId, nil
	}

	// Cache is used, but no data is ready.
	var validator *Validator
	validator, err = km.validateToken(tokenString)
	if err != nil {
		return validator.userId, validator.sessionId, err
	}

	cd = &CachedData{
		UserId:         validator.userId,
		SessionId:      validator.sessionId,
		ExpirationTime: validator.expTime,
	}
	err = km.cache.AddRecord(tokenString, cd)
	if err != nil {
		return validator.userId, validator.sessionId, err
	}

	return validator.userId, validator.sessionId, nil
}

func (km *KeyMaker) validateToken(tokenString string) (validator *Validator, err error) {
	validator = NewValidator(km.signingMethodName, km.publicKey)

	var token *jwt.Token
	token, err = jwt.Parse(tokenString, validator.KeyFn)
	if err != nil {
		return validator, err
	}

	if !token.Valid {
		return validator, errors.New(ErrTokenIsNotValid)
	}

	return validator, nil
}
func (km *KeyMaker) getCachedData(tokenString string) (cd *CachedData) {
	var err error
	cd, err = km.cache.GetRecord(tokenString)
	if err != nil {
		return nil
	}
	if cd == nil {
		return nil
	}
	if !cd.IsGood() {
		return nil
	}

	return cd
}
