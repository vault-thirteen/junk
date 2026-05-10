package km

import (
	"fmt"
	"testing"
	"time"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_KeyMaker_ValidateToken(t *testing.T) {
	const (
		Common_SigningMethodName  = TokenAlg_RS512
		Common_PathPrefix         = `..\..\..\`
		Common_PrivateKeyFilePath = Common_PathPrefix + "cert\\JWT\\jwtPrivateKey.pem"
		Common_PublicKeyFilePath  = Common_PathPrefix + "cert\\JWT\\jwtPublicKey.pem"
	)

	userId := 123
	sessionId := 456
	var a, b int

	aTest := tester.New(t)
	n := 20_000
	resumeFN := func(testNumber int, dur time.Duration) {
		msg := fmt.Sprintf("Test #%v. Duration: %v ms.", testNumber, dur.Milliseconds())
		fmt.Println(msg)
	}

	// Test #1. Cache is enabled.
	{
		settings := &KeyMakerSettings{
			SigningMethodName:  Common_SigningMethodName,
			PrivateKeyFilePath: Common_PrivateKeyFilePath,
			PublicKeyFilePath:  Common_PublicKeyFilePath,
			IsCacheEnabled:     true,
			CacheSizeLimit:     1024,
			CacheRecordTtl:     60,
		}

		keyMaker, err := New(settings)
		aTest.MustBeNoError(err)

		var tokenString string
		tokenString, err = keyMaker.MakeJWToken(userId, sessionId, time.Now().Add(time.Minute))
		aTest.MustBeNoError(err)

		var t1 = time.Now()
		for i := 0; i < n; i++ {
			a, b, err = keyMaker.ValidateToken(tokenString)
			aTest.MustBeNoError(err)
			aTest.MustBeEqual(a, userId)
			aTest.MustBeEqual(b, sessionId)
		}
		var dur = time.Now().Sub(t1)
		resumeFN(1, dur)
	}

	// Test #2. Cache is disabled.
	{
		settings := &KeyMakerSettings{
			SigningMethodName:  Common_SigningMethodName,
			PrivateKeyFilePath: Common_PrivateKeyFilePath,
			PublicKeyFilePath:  Common_PublicKeyFilePath,
			IsCacheEnabled:     false,
		}

		keyMaker, err := New(settings)
		aTest.MustBeNoError(err)

		var tokenString string
		tokenString, err = keyMaker.MakeJWToken(userId, sessionId, time.Now().Add(time.Minute))
		aTest.MustBeNoError(err)

		var t1 = time.Now()
		for i := 0; i < n; i++ {
			a, b, err = keyMaker.ValidateToken(tokenString)
			aTest.MustBeNoError(err)
			aTest.MustBeEqual(a, userId)
			aTest.MustBeEqual(b, sessionId)
		}
		var dur = time.Now().Sub(t1)
		resumeFN(2, dur)
	}
}
