package assert

import "testing"

func Eq(t *testing.T, expected any, actual any) {
	if expected != actual {
		t.Errorf("Expected: %v\nGot: %v", expected, actual)
	}
}
