// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
