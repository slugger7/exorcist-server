package assert

import (
	"errors"
	"testing"

	errs "github.com/slugger7/exorcist/internal/errors"
)

func ErrorNil(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected an error but it was nil")
	}
}

func ErrorMessage(t *testing.T, err error, message string) {
	var e errs.IError
	if errors.As(err, &e) {
		if e.Message() != message {
			t.Errorf("Expected error message: %v\nGot error message: %v", message, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got: %v", err.Error())
	}
}
