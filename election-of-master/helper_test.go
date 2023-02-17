// helper_test.go.

package eom

import (
	"testing"

	"github.com/vault-thirteen/tester"
)

func Test_checkServiceName(t *testing.T) {

	var aTest *tester.Test
	var err error
	var serviceName string

	aTest = tester.New(t)

	// Test #1. Good Name.
	serviceName = "abcXYZ_123"
	err = checkServiceName(serviceName)
	aTest.MustBeNoError(err)

	// Test #2. Bad Name
	serviceName = "ы"
	err = checkServiceName(serviceName)
	aTest.MustBeAnError(err)

	// Test #3. Bad Name
	serviceName = "%"
	err = checkServiceName(serviceName)
	aTest.MustBeAnError(err)
}
