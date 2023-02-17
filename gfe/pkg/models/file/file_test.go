package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	// Arrange.
	resultExpected := ID("x")

	// Act.
	resultActual := NewID("x")

	// Assert.
	assert.Equal(t, resultExpected, resultActual)
}
