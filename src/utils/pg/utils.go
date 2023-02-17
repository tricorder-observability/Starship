package pg

import "strings"

// pgPath read string slice, return string with postgres support format, e.g. ->'key'->>'subKey'
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
