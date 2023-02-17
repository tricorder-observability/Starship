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
