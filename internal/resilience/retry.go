package resilience

import (
	"time"
)

// Retry retries a function multiple times
func Retry(attempts int, sleep time.Duration, fn func() error) error {

	var err error

	for i := 0; i < attempts; i++ {

		err = fn()

		if err == nil {
			return nil
		}

		time.Sleep(sleep)
	}

	return err
}
