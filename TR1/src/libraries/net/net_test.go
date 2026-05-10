package net

import (
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_SplitUrlPath(t *testing.T) {
	aTest := tester.New(t)

	aTest.MustBeEqual(SplitUrlPath(""), []string{})
	aTest.MustBeEqual(SplitUrlPath("/"), []string{})
	aTest.MustBeEqual(SplitUrlPath("//"), []string{})
	aTest.MustBeEqual(SplitUrlPath("a/"), []string{"a"})
	aTest.MustBeEqual(SplitUrlPath("/a"), []string{"a"})
	aTest.MustBeEqual(SplitUrlPath("a/b"), []string{"a", "b"})
	aTest.MustBeEqual(SplitUrlPath("a///b"), []string{"a", "b"})
}
