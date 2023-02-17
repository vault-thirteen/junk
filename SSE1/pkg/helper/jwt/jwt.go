package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vault-thirteen/auxie/random"
	randomHelper "github.com/vault-thirteen/junk/SSE1/pkg/helper/random"
)

const (
	ErrNullPointer     = "Null Pointer"
	ErrTypeCastFailure = "Type Cast Failure"
)

const (
	TokenKeyBytesCount = 32
	TokenIssuer        = "Vault Thirteen"
	TokenAudience      = "Humanity"

	// JSON Web Token Claim Names.
	// As per Documentation at:
	// https://openid.net/specs/draft-jones-json-web-token-07.html
	TokenClaimExpirationTime = "exp"
	TokenClaimNotBefore      = "nbf"
	TokenClaimIssuedAt       = "iat"
	TokenClaimIssuer         = "iss"
	TokenClaimAudience       = "aud"
	TokenClaimPrincipal      = "prn"
	TokenClaimJWTID          = "jti"
	TokenClaimType           = "typ"
	TokenClaimAlgorithm      = "alg"

	// Custom Claims.
	TokenClaimSessionId  = "sid"
	TokenClaimMarkerHash = "mrh"
)

// Auxiliary Data used with a JSON Web Token in this Project.
type TokenData struct {

	// A random Key.
	TokenKey    []byte
	TokenKeyStr string

	// Token as Text.
	TokenString string

	// Internal Data of a Token...

	// 1. A unique Marker with its SHA-256 Hash Sum.
	UniqueMarker     string
	UniqueMarkerHash string

	// 2. Id of a Session.
	SessionId uint
}

// Creates a set of auxiliary Data used with a JSON Web Token which includes:
//   - A unique random Marker,
//   - A SHA-256 Hashsum of the above Marker,
//   - A random Token Key represented as an Array (Slice) of Bytes,
//   - A hexadecimal String of the above Key.
func PrepareDataForToken() (td *TokenData, err error) {

	// Marker.
	var uniqueMarker string
	uniqueMarker, err = randomHelper.CreateUniqueMarker()
	if err != nil {
		return
	}
	var markerBytes []byte
	markerBytes, err = hex.DecodeString(uniqueMarker)
	if err != nil {
		return
	}
	markerHash := sha256.Sum256(markerBytes)
	td = &TokenData{
		UniqueMarker:     uniqueMarker,
		UniqueMarkerHash: strings.ToUpper(hex.EncodeToString((markerHash)[:])),
	}

	// Key.
	td.TokenKey, err = createTokenKey()
	if err != nil {
		return
	}
	td.TokenKeyStr = strings.ToUpper(hex.EncodeToString(td.TokenKey))
	return
}

// Creates a JSON Web Token.
func CreateJWToken(
	sessionId uint,
	markerHash string,
	tokenLifeTimeSec uint,
) (token *jwt.Token, err error) {
	var timeNowTS = time.Now().Unix()
	var claims = jwt.MapClaims{
		TokenClaimIssuer:   TokenIssuer,
		TokenClaimAudience: TokenAudience,
		//
		TokenClaimIssuedAt:       timeNowTS,
		TokenClaimNotBefore:      timeNowTS,
		TokenClaimExpirationTime: timeNowTS + int64(tokenLifeTimeSec),
		//
		TokenClaimSessionId:  sessionId,
		TokenClaimMarkerHash: markerHash,
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return
}

// Creates a random Key for a J.W. Token.
func createTokenKey() (tokenKey []byte, err error) {
	return random.GenerateRandomBytesA1(TokenKeyBytesCount)
}

// Verifies a J.W. Token.
func VerifyJWToken(
	tokenString string,
	keyFunc jwt.Keyfunc,
) (token *jwt.Token, err error) {
	token, err = jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return
	}
	return
}

// Reads some Data from a J.W. Token including:
//   - Marker Hash Sum,
//   - Session Id.
func GetTokenData(
	token *jwt.Token,
) (sessionId uint, markerHash string, err error) {
	if token == nil {
		err = errors.New(ErrNullPointer)
		return
	}
	var claims jwt.MapClaims
	claims, err = GetTokenClaims(token)
	if err != nil {
		return
	}
	markerHash, err = GetTokenClaimMarkerHash(claims)
	if err != nil {
		return
	}
	sessionId, err = GetTokenClaimSessionId(claims)
	if err != nil {
		return
	}
	return
}

// Reads the 'claims' Object from the J.W. Token.
func GetTokenClaims(
	token *jwt.Token,
) (claims jwt.MapClaims, err error) {
	if token == nil {
		err = errors.New(ErrNullPointer)
		return
	}
	var ok bool
	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New(ErrTypeCastFailure)
		return
	}
	return
}

// Reads the Marker Hash Sum from the 'claims' Object.
func GetTokenClaimMarkerHash(
	claims jwt.MapClaims,
) (markerHash string, err error) {
	var markerHashIfc interface{}
	var ok bool
	markerHashIfc, ok = claims[TokenClaimMarkerHash]
	if !ok {
		err = errors.New(ErrTypeCastFailure)
		return
	}
	markerHash, ok = markerHashIfc.(string)
	if !ok {
		err = errors.New(ErrTypeCastFailure)
		return
	}
	return
}

// Reads the Session Id from the 'claims' Object.
func GetTokenClaimSessionId(
	claims jwt.MapClaims,
) (sessionId uint, err error) {
	var sessionIdIfc interface{}
	var ok bool
	sessionIdIfc, ok = claims[TokenClaimSessionId]
	if !ok {
		err = errors.New(ErrTypeCastFailure)
		return
	}
	var sessionIdFloat float64
	sessionIdFloat, ok = sessionIdIfc.(float64)
	if !ok {
		err = errors.New(ErrTypeCastFailure)
		return
	}
	sessionId = uint(sessionIdFloat)
	return
}
