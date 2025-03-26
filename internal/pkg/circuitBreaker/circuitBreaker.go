package circuitbreaker

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker/v2"
)

func NewCircuitBreaker[T any](name string) *gobreaker.CircuitBreaker[T] {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 5,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit breaker state changed from %s to %s\n", from, to)
		},
	}

	return gobreaker.NewCircuitBreaker[T](settings)
}
