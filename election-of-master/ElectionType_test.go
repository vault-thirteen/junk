// ElectionType_test.go.

package eom

import (
	"testing"

	"github.com/vault-thirteen/tester"
)

func Test_IsValid(t *testing.T) {

	var aTest *tester.Test
	var et ElectionType

	aTest = tester.New(t)

	// Test #1. Junk.
	et = ElectionType(0)
	aTest.MustBeEqual(et.IsValid(), false)

	// Test #1. Normal Type.
	et = ElectionType(ElectionTypeSingleMaster)
	aTest.MustBeEqual(et.IsValid(), true)
}
