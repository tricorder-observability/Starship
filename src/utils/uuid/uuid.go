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

package uuid

import (
	"strings"

	"github.com/google/uuid"
)

// New returns a UUID.
func New() string {
	return uuid.New().String()
}

// Returns a UUID with the provided separator.
func NewWithSeparator(separator string) string {
	uuid := New()
	const defaultSeparator = "-"
	return strings.ReplaceAll(uuid, defaultSeparator, separator)
}

// Returns a UUID with underscore as the separator.
func NewWithUnderscoreSeparator() string {
	const underscore = "_"
	return NewWithSeparator(underscore)
}
