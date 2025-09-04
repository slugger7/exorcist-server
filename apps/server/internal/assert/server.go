package assert

import (
	"testing"
)

// Deprecated: use testify assert
func StatusCode(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("Expected status code: %v\nGot status code: %v", expectedCode, actualCode)
	}
}

// Deprecated: use testify assert
func Body(t *testing.T, expectedBody, actualBody string) {
	if expectedBody != actualBody {
		t.Errorf("Expected body: %v\nGot body: %v", expectedBody, actualBody)
	}
}
