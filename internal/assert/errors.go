package assert

import (
	"errors"
	"testing"

	errs "github.com/slugger7/exorcist/internal/errors"
)

func ErrorNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected nil error but was: %v", err.Error())
	}
}

func ErrorNotNil(t *testing.T, err error) {
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

func Error(t *testing.T, expectedErr, err error) {
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedErr.Error(), err.Error())
	}
}
