// Copyright (C) 2023  tricorder-observability
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

package bytes

import "bytes"

// TrimAfter returns a byte slice with the first appearance of `c` and all its trailing bytes removed.
func TrimAfter(s []byte, sep byte) []byte {
	pos := bytes.IndexByte(s, sep)
	if pos == -1 {
		return s
	}
	return s[:pos]
}

// TrimC returns a byte slice with the first appearance of `\x00` (null character) and all its trailing bytes removed.
func TrimC(s []byte) []byte {
	return TrimAfter(s, '\x00')
}
