package models

import (
	"fmt"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewFormatString(t *testing.T) {
	aTest := tester.New(t)
	var fs *FormatString
	var err error

	type TestCase struct {
		S                    string
		IsErrorExpected      bool
		ExpectedType         string
		ExpectedHasM         bool
		ExpectedHasR         bool
		ExpectedHasT         bool
		ExpectedHasU         bool
		ExpectedPlaceholders []Placeholder
	}
	var tests = []TestCase{
		{
			S:               "Hello World",
			IsErrorExpected: false,
			ExpectedType:    "",
		},
		{
			S:               "Aaa {M} bbb {R} ccc {T} ddd {U} eee.",
			IsErrorExpected: false,
			ExpectedType:    "MRTU",
			ExpectedHasM:    true,
			ExpectedHasR:    true,
			ExpectedHasT:    true,
			ExpectedHasU:    true,
		},
		{
			S:               "Aaa {T} bbb {U} ccc {M} ddd.",
			IsErrorExpected: false,
			ExpectedType:    "TUM",
			ExpectedHasM:    true,
			ExpectedHasT:    true,
			ExpectedHasU:    true,
			ExpectedPlaceholders: []Placeholder{
				{Type: PlaceholderTypeT, Pos: 4},
				{Type: PlaceholderTypeU, Pos: 12},
				{Type: PlaceholderTypeM, Pos: 20},
			},
		},
		{
			S:               "{R} ... {T}",
			IsErrorExpected: false,
			ExpectedType:    "RT",
			ExpectedHasR:    true,
			ExpectedHasT:    true,
		},
		{
			S:               "{R}",
			IsErrorExpected: false,
			ExpectedType:    "R",
			ExpectedHasR:    true,
		},
		{
			S:               "Duplicate M: {M} and {M}.",
			IsErrorExpected: true,
		},
		{
			S:               "Duplicate R: {R} and {R}.",
			IsErrorExpected: true,
		},
		{
			S:               "Duplicate T: {T} and {T}.",
			IsErrorExpected: true,
		},
		{
			S:               "Duplicate U: {U} and {U}.",
			IsErrorExpected: true,
		},
		{
			S:               "Not a format: {X}, {Y}, {Z}.",
			IsErrorExpected: false,
			ExpectedType:    "",
		},
	}

	for i, test := range tests {
		fmt.Print(fmt.Sprintf("[%d]", i+1))
		fs, err = NewFormatString(test.S)
		if test.IsErrorExpected {
			aTest.MustBeAnError(err)
		} else {
			aTest.MustBeNoError(err)

			aTest.MustBeEqual(fs.String(), test.S)
			aTest.MustBeEqual(fs.Type(), test.ExpectedType)
			aTest.MustBeEqual(fs.HasM(), test.ExpectedHasM)
			aTest.MustBeEqual(fs.HasR(), test.ExpectedHasR)
			aTest.MustBeEqual(fs.HasT(), test.ExpectedHasT)
			aTest.MustBeEqual(fs.HasU(), test.ExpectedHasU)

			if test.ExpectedPlaceholders != nil {
				aTest.MustBeEqual(fs.Placeholders(), test.ExpectedPlaceholders)
			}
		}
	}
	fmt.Println()
}
