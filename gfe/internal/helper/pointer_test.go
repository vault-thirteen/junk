package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStringPointer(t *testing.T) {
	// Act.
	ptr := NewStringPointer("x")

	// Assert.
	assert.NotEqual(t, nil, ptr)
	assert.Equal(t, "x", *ptr)
}

func TestNewIntPointer(t *testing.T) {
	// Act.
	ptr := NewIntPointer(123)

	// Assert.
	assert.NotEqual(t, nil, ptr)
	assert.Equal(t, 123, *ptr)
}
