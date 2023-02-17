package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcatenateEnvVarPrefixes(t *testing.T) {
	// Act & assert.
	assert.Equal(t, "a_b", ConcatenateEnvVarPrefixes("a", "b"))
}
