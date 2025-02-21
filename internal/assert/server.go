package assert

import (
	"testing"
)

func StatusCode(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("Expected status code: %v\nGot status code: %v", expectedCode, actualCode)
	}
}

func Body(t *testing.T, expectedBody, actualBody string) {
	if expectedBody != actualBody {
		t.Errorf("Expected body: %v\nGot body: %v", expectedBody, actualBody)
	}
}
