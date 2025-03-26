package circuitbreaker_test

import (
	circuitbreaker "main-api/internal/pkg/circuitBreaker"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker(t *testing.T) {
	name := "test-circuit"

	cb := circuitbreaker.NewCircuitBreaker[any](name)
	assert.NotNil(t, cb, "Circuit breaker should not be nil")
	assert.Equal(t, name, cb.Name(), "Circuit breaker name should match")
}
