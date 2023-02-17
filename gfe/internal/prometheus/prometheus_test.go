package prometheus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPrometheus(t *testing.T) {
	// Act.
	p, err := NewPrometheus()

	// Assert.
	assert.NotEqual(t, nil, p)
	assert.Equal(t, err, nil)
}

func TestPrometheus_GetMetrics(t *testing.T) {
	// Arrange.
	p, err := NewPrometheus()
	assert.Equal(t, err, nil)

	// Act.
	m := p.GetMetrics()

	// Assert.
	assert.NotEqual(t, nil, m)
}
