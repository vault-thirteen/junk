package dk

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewDKey(t *testing.T) {
	aTest := tester.New(t)

	var key *DKey
	var err error

	// Test #1. Wrong key size.
	key, err = NewDKey(0)
	aTest.MustBeAnError(err)

	// Test #2. Normal key size.
	key, err = NewDKey(4)
	aTest.MustBeNoError(err)
	aTest.MustBeDifferent(key, (*DKey)(nil))
}

func Test_GetBytes(t *testing.T) {
	aTest := tester.New(t)

	var key *DKey
	var err error
	var bytes []byte

	key, err = NewDKey(4)
	aTest.MustBeNoError(err)

	// Test #1. First reading.
	bytes = key.GetBytes()
	aTest.MustBeDifferent(bytes, ([]byte)(nil))
	aTest.MustBeEqual(len(bytes), 4)

	// Test #2. Second reading.
	bytes = key.GetBytes()
	aTest.MustBeEqual(bytes, ([]byte)(nil))
	aTest.MustBeEqual(len(bytes), 0)
}

func Test_GetString(t *testing.T) {
	aTest := tester.New(t)

	var key *DKey
	var err error
	var str string

	key, err = NewDKey(4)
	aTest.MustBeNoError(err)

	// Test #1. First reading.
	str = key.GetString()
	aTest.MustBeDifferent(str, "")
	aTest.MustBeEqual(len(str), 4*2)

	// Test #2. Second reading.
	str = key.GetString()
	aTest.MustBeEqual(str, "")
	aTest.MustBeEqual(len(str), 0)
}

func Test_CheckBytes(t *testing.T) {
	aTest := tester.New(t)

	var key *DKey
	var err error
	var x []byte

	key, err = NewDKey(32)
	aTest.MustBeNoError(err)

	// Test #1. Negative check.
	x = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	aTest.MustBeEqual(key.CheckBytes(x), false)

	// Test #2. Positive check.
	x = key.bytes
	aTest.MustBeEqual(key.CheckBytes(x), true)
}

func Test_CheckString(t *testing.T) {
	aTest := tester.New(t)

	var key *DKey
	var err error
	var x string

	key, err = NewDKey(32)
	aTest.MustBeNoError(err)

	// Test #1. Negative check.
	x = "0000000000000000000000000000000000000000000000000000000000000000"
	aTest.MustBeEqual(key.CheckString(x), false)

	// Test #2. Positive check.
	x = key.str
	aTest.MustBeEqual(key.CheckString(x), true)
}
