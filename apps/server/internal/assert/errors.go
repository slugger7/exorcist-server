package assert

import (
	"errors"
	"testing"

	errs "github.com/slugger7/exorcist/internal/errors"
)

// Deprecated: use testify assert
func ErrorNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected nil error but was: %v", err.Error())
	}
}

// Deprecated: use testify assert
func ErrorNotNil(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected an error but it was nil")
	}
}

// Deprecated: use testify assert
func ErrorMessage(t *testing.T, message string, err error) {
	var e errs.IError
	if errors.As(err, &e) {
		if e.Message() != message {
			t.Errorf("Expected error message: %v\nGot error message: %v", message, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got: %v", err.Error())
	}
}

// Deprecated: use testify assert
func Error(t *testing.T, expectedErr, err error) {
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedErr.Error(), err.Error())
	}
}
