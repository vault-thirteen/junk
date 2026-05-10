package km

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

//TODO:RSA-PSS.
// Unfortunately, Golang can not parse RSA-PSS (RSASSA-PSS) keys.
// See this bug report: https://github.com/golang/go/issues/48314
// What a shame again. Guys, why is your language so poor ?
// x509: PKCS#8 wrapping contained private key with unknown algorithm: 1.2.840.113549.1.1.10
//var anyKey any
//anyKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
//if err != nil {
//	return nil, err
//}

const (
	ErrRsaPssGoLanguageWhatAShame = "this stupid Go language does not support parsing RSA-PSS keys, what a shame"
)

// This function is bad, but Golang is ever worse.
func getPrivateKey(privateKeyFilePath string, signingMethodName string) (pk *rsa.PrivateKey, err error) {
	var keyFileData []byte
	keyFileData, err = os.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyFileData)
	var anyKey any
	var ok bool

	switch signingMethodName {
	case TokenAlg_PS512:
		return nil, errors.New(ErrRsaPssGoLanguageWhatAShame)
	case TokenAlg_RS512:
		{
			anyKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}

			pk, ok = anyKey.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New(ErrTypeCast)
			}

			return pk, nil
		}
	default:
		return nil, errors.New(ErrSigningMethodIsNotSupported)
	}
}

// This function is bad, but Golang is ever worse.
func getPublicKey(publicKeyFilePath string, signingMethodName string) (pk *rsa.PublicKey, err error) {
	var keyFileData []byte
	keyFileData, err = os.ReadFile(publicKeyFilePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyFileData)
	var anyKey any
	var ok bool

	switch signingMethodName {
	case TokenAlg_PS512:
		return nil, errors.New(ErrRsaPssGoLanguageWhatAShame)
	case TokenAlg_RS512:
		{
			anyKey, err = x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}

			pk, ok = anyKey.(*rsa.PublicKey)
			if !ok {
				return nil, errors.New(ErrTypeCast)
			}

			return pk, nil
		}
	default:
		return nil, errors.New(ErrSigningMethodIsNotSupported)
	}
}
