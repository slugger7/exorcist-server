package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func Test_CreateLibrary_InvalidBody(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	r.POST("/", s.server.CreateLibrary)

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

func Test_CreateLibrary_NoNameSpecified_ShouldThrowError(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	r.POST("/", s.server.CreateLibrary)

	req, err := http.NewRequest("POST", "/", body(`{"name": ""}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"no value for name"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_CreateLibrary_ErrorByService(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	expectedErrorMessage := "expected error message"
	s.mockService.Library.MockError[0] = errors.New(expectedErrorMessage)
	r.POST("/", s.server.CreateLibrary)

	expectedName := "expectedLibraryName"
	req, err := http.NewRequest("POST", "/", body(`{"name":"%v"}`, expectedName))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"could not create new library"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_CreateLibrary_Success(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	expectedId, _ := uuid.NewRandom()
	expectedLibraryName := "some expected library name"
	s.mockService.Library.MockModel[0] = &model.Library{
		ID:   expectedId,
		Name: expectedLibraryName,
	}

	r.POST("/", s.server.CreateLibrary)

	req, err := http.NewRequest("POST", "/", body(`{"name":"%v"}`, expectedLibraryName))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusCreated
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"id":"%v"}`, expectedId.String())
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_GetLibraries_ServiceReturnsError(t *testing.T) {
	r := setupEngine()
	s := setupServer()
	expectedError := errors.New("expected error")
	s.mockService.Library.MockError[0] = expectedError

	r.GET("/", s.server.GetLibraries)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusInternalServerError
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"could not fetch libraries"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_LibraryAction_WithInvalidId(t *testing.T) {
	s := setupServer()

	invalidId := "not-a-uuid"
	rr := s.withGetEndpoint(s.server.LibraryAction, ":id/*action").
		withGetRequest(nil, fmt.Sprintf("%v/action", invalidId)).
		exec()

	expectedStatus := http.StatusBadRequest
	if rr.Code != expectedStatus {
		t.Errorf("Exected status code: %v\nGot status code: %v", expectedStatus, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"error":"%v"}`, fmt.Sprintf(ErrIdParse, invalidId))
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body: %v\nGot body: %v", expectedBody, rr.Body.String())
	}
}

func Test_LibraryAction_WithServiceReturningError(t *testing.T) {
	s := setupServer()

	s.mockService.Library.MockError[0] = fmt.Errorf("error")

	id, _ := uuid.NewRandom()
	action := "action"

	rr := s.withGetEndpoint(s.server.LibraryAction, ":id/*action").
		withGetRequest(nil, fmt.Sprintf("%v/%v", id, action)).
		exec()

	expectedStatus := http.StatusInternalServerError
	if rr.Code != expectedStatus {
		t.Errorf("Exected status code: %v\nGot status code: %v", expectedStatus, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"error":"%v"}`, fmt.Sprintf(ErrLibraryAction, "/"+action, id))
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body: %v\nGot body: %v", expectedBody, rr.Body.String())
	}
}
