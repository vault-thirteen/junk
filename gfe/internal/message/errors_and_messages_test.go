package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeMessageWithPrefix(t *testing.T) {
	// Act & assert.
	assert.Equal(t, "a b", ComposeMessageWithPrefix("a", "b"))
}
