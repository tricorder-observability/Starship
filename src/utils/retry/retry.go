// Package retry provides API to retry a function if it fails.
package retry

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// ExpBackOffWithLimit retries the input function up to 10 times with the default exponential backoff.
func ExpBackOffWithLimit(fn func() error) error {
	backOff := backoff.NewExponentialBackOff()
	backOff.MaxInterval = 10 * time.Second

	const retryLimit = 10
	defaultBackoff := backoff.WithMaxRetries(backOff, retryLimit /*max_retry*/)
	err := backoff.Retry(fn, defaultBackoff)
	if err != nil {
		return fmt.Errorf("failed after retrying %d times, last error: %v", retryLimit, err)
	}
	return nil
}
