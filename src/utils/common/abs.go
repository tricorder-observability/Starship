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

package common

// TODO(yzhao): Use generic
// Golang does not have Abs for integers. See https://stackoverflow.com/a/57649529
func AbsInt8(v int8) int {
	if v < 0 {
		return int(-v)
	}
	return int(v)
}

func AbsUint8s(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}
