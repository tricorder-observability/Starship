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

package pg

import "strings"

// pgPath read string slice, return string with postgres support format, e.g. ->'key'->>'subKey'.
func pgPath(paths []string) string {
	if len(paths) == 0 {
		return "data->'metadata'->>'uid'"
	}
	sb := strings.Builder{}
	for i := range paths {
		if i != len(paths)-1 {
			sb.WriteString("->'" + paths[i] + "'")
		} else {
			sb.WriteString("->>'" + paths[i] + "'")
		}
	}
	return "data" + sb.String()
}
