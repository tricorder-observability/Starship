package common

import "strings"

// StrTrimPrefix returns a string with the leading l chars removed.
func StrTrimPrefix(s string, l int) string {
	return s[l:]
}

// StrTrimSuffix returns a string with the trailing l chars removed.
func StrTrimSuffix(s string, l int) string {
	return s[:len(s)-l]
}

// StrTrimAfter returns a string with the first appearance of `c` and all its trailing bytes removed.
func StrTrimAfter(s, sep string) string {
	pos := strings.Index(s, sep)
	return s[:pos]
}
