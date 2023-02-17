package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeUniqueRandomString(t *testing.T) {
	// Act.
	s := MakeUniqueRandomString()

	// Assert.
	assert.Equal(t, true, len(s) > 0)
}
