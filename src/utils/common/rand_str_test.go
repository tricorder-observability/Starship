package common

import (
	"testing"
)

func TestRandStr(t *testing.T) {
	str1 := RandStr(10)
	str2 := RandStr(10)
	if str1 == str2 {
		t.Errorf("Random string generated from 2 RandStr() should be different, got '%s'", str1)
	}
}
