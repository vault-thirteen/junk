package main

import (
	"errors"
	"flag"
	"fmt"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"log"

	"github.com/vault-thirteen/SimpleBB/pkg/ACM/km"
)

const (
	ErrUserIdIsNotSet    = "user ID is not set"
	ErrSessionIdIsNotSet = "session ID is not set"
	ErrKeyFileIsNotSet   = "key file is not set"
)

func main() {
	userId, sessionId, privateKeyFilePath, publicKeyFilePath, signingMethod, err := receiveArguments()
	mustBeNoError(err)

	var keyMaker *km.KeyMaker
	keyMaker, err = km.New(signingMethod, privateKeyFilePath, publicKeyFilePath)
	mustBeNoError(err)

	var ts simple.WebTokenString
	ts, err = keyMaker.MakeJWToken(userId, sessionId)
	mustBeNoError(err)

	fmt.Println(fmt.Sprintf("Token string: %v.", ts))

	userId, sessionId, err = keyMaker.ValidateToken(ts)
	mustBeNoError(err)

	fmt.Println(fmt.Sprintf("userId=%v, sessionId=%v, signingMethod=%v.", userId, sessionId, signingMethod))
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func receiveArguments() (userId cmb.Id, sessionId cmb.Id, privateKeyFilePath simple.Path, publicKeyFilePath simple.Path, signingMethod string, err error) {
	var userIdInt int
	flag.IntVar(&userIdInt, "uid", 0, "user ID")
	var sessionIdInt int
	flag.IntVar(&sessionIdInt, "sid", 0, "session ID")
	var privateKeyFilePathStr string
	flag.StringVar(&privateKeyFilePathStr, "private_key", "", "path to private key file using PEM format")
	var publicKeyFilePathStr string
	flag.StringVar(&publicKeyFilePathStr, "public_key", "", "path to public key file using PEM format")
	flag.StringVar(&signingMethod, "method", "", "signing method")
	flag.Parse()

	userId = cmb.Id(userIdInt)
	sessionId = cmb.Id(sessionIdInt)
	privateKeyFilePath = simple.Path(privateKeyFilePathStr)
	publicKeyFilePath = simple.Path(publicKeyFilePathStr)

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
