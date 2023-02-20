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
