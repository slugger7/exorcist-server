package assert

import "testing"

// Deprecated: use testify assert
func Eq(t *testing.T, expected any, actual any) {
	if expected != actual {
		t.Errorf("\nExpected: %v\nGot     : %v", expected, actual)
	}
}
