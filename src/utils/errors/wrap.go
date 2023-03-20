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

package errors

import (
	"errors"
	"fmt"
)

// Wrap returns a new error object with the context enclosed in its message.
func Wrap(context, failure string, err error) error {
	return fmt.Errorf("while %s, failed to %s, error: %v", context, failure, err)
}

// New returns a new error object with the context.
func New(context, failure string) error {
	return fmt.Errorf("%s, failed to %s", context, failure)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
