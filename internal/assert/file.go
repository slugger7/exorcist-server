package assert

import (
	"os"
	"testing"
)

func FileExists(t *testing.T, file string) {
	if _, err := os.Stat(file); err != nil {
		t.Errorf("file did not exist %v: %v", file, err.Error())
	}
}
