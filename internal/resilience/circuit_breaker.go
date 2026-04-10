package resilience

import (
	"errors"
	"sync"
	"time"
)

type CircuitBreaker struct {
	failures  int
	threshold int
	open      bool
	mutex     sync.Mutex
	lastFail  time.Time
	timeout   time.Duration
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		timeout:   timeout,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {

	cb.mutex.Lock()

	// if open, check timeout
	if cb.open {
		if time.Since(cb.lastFail) > cb.timeout {
			cb.open = false // half-open
		} else {
			cb.mutex.Unlock()
			return errors.New("circuit breaker open")
		}
	}

	cb.mutex.Unlock()

	err := fn()

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFail = time.Now()

		if cb.failures >= cb.threshold {
			cb.open = true
		}
		return err
	}

	// reset on success
	cb.failures = 0
	return nil
}
