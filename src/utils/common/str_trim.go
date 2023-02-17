package common

// StrTrimPrefix returns a string with the leading l chars removed.
func StrTrimPrefix(s string, l int) string {
	return s[l:]
}

// StrTrimSuffix returns a string with the trailing l chars removed.
func StrTrimSuffix(s string, l int) string {
	return s[:len(s)-l]
}
