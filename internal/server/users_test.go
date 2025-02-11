package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func Test_Create_InvalidBody(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	r.POST("/", s.server.CreateUser)

	req, err := http.NewRequest("POST", "/", body(`{invalid json}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"could not read body of request"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Create_ServiceReturnsError(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	expectedErrorMessage := "expected error"
	s.mockService.UserService.MockErrors[0] = errors.New(expectedErrorMessage)
	r.POST("/", s.server.CreateUser)

	req, err := http.NewRequest("POST", "/", body(`{"username":"someUsername","password":"somePassword"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"error":"%s"}`, expectedErrorMessage)
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_Create_Success(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	expectedModel := &model.User{
		Username: "expecedUsername",
		Password: "",
	}
	s.mockService.UserService.MockModel[0] = expectedModel

	r.POST("/", s.server.CreateUser)

	req, err := http.NewRequest("POST", "/", body(fmt.Sprintf(`{"username":"%s","password":"somePassword"}`, expectedModel.Username)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusCreated
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"ID":"00000000-0000-0000-0000-000000000000","Username":"%s","Password":"","Active":false,"Created":"0001-01-01T00:00:00Z","Modified":"0001-01-01T00:00:00Z"}`, expectedModel.Username)
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}
