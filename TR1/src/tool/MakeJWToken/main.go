package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/vault-thirteen/TR1/src/libraries/km"
)

const (
	ErrUserIdIsNotSet    = "user ID is not set"
	ErrSessionIdIsNotSet = "session ID is not set"
	ErrKeyFileIsNotSet   = "key file is not set"
)

func main() {
	userId, sessionId, privateKeyFilePath, publicKeyFilePath, signingMethod, err := receiveArguments()
	mustBeNoError(err)

	kms := &km.KeyMakerSettings{
		SigningMethodName:  signingMethod,
		PrivateKeyFilePath: privateKeyFilePath,
		PublicKeyFilePath:  publicKeyFilePath,
		IsCacheEnabled:     true,
		CacheSizeLimit:     1024,
		CacheRecordTtl:     60,
	}

	var keyMaker *km.KeyMaker
	keyMaker, err = km.New(kms)
	mustBeNoError(err)

	sessionMaxDurationSec := 5 * 60
	expirationTime := time.Now().Add(time.Duration(sessionMaxDurationSec) * time.Second)

	var ts string
	ts, err = keyMaker.MakeJWToken(userId, sessionId, expirationTime)
	mustBeNoError(err)

	fmt.Println(fmt.Sprintf("Token string: %v.", ts))

	userId, sessionId, err = keyMaker.ValidateToken(ts)
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func receiveArguments() (userId int, sessionId int, privateKeyFilePath string, publicKeyFilePath string, signingMethod string, err error) {
	flag.IntVar(&userId, "uid", 0, "user ID")
	flag.IntVar(&sessionId, "sid", 0, "session ID")
	flag.StringVar(&privateKeyFilePath, "private_key", "", "path to private key file using PEM format")
	flag.StringVar(&publicKeyFilePath, "public_key", "", "path to public key file using PEM format")
	flag.StringVar(&signingMethod, "method", "", "signing method")
	flag.Parse()

	if userId == 0 {
		return 0, 0, "", "", "", errors.New(ErrUserIdIsNotSet)
	}

	if sessionId == 0 {
		return 0, 0, "", "", "", errors.New(ErrSessionIdIsNotSet)
	}

	if len(privateKeyFilePath) == 0 {
		return 0, 0, "", "", "", errors.New(ErrKeyFileIsNotSet)
	}

	if len(publicKeyFilePath) == 0 {
		return 0, 0, "", "", "", errors.New(ErrKeyFileIsNotSet)
	}

	return userId, sessionId, privateKeyFilePath, publicKeyFilePath, signingMethod, nil
}
