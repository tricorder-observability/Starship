package errors

import "fmt"

// Wrap returns a new error object with the context enclosed in its message.
func Wrap(context, failure string, err error) error {
	return fmt.Errorf("%s, failed to %s, error: %v", context, failure, err)
}
