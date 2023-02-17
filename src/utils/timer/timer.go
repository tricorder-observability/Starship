package timer

import "time"

type Timer struct {
	startTime time.Time
}

// New returns a timer that use now as the start time.
func New() *Timer {
	timer := new(Timer)
	timer.startTime = time.Now()
	return timer
}

// Stop returns the elapsed time since the start.
func (t *Timer) Get() time.Duration {
	return time.Since(t.startTime)
}
