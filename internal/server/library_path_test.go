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

func Test_CreateLibraryPath_InvalidBody(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	r.POST("/", s.server.CreateLibraryPath)

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

func Test_CreateLibraryPath_NoPathSpecified_ShouldThrowError(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	r.POST("/", s.server.CreateLibraryPath)

	req, err := http.NewRequest("POST", "/", body(`{"path": ""}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusBadRequest
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := `{"error":"no path in body"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

// func Test_CreateLibraryPath_ErrorWhileFetchingLibraryById(t *testing.T) {
// 	r := setupEngine()
// 	s := setupServer()
// }

func Test_CreateLibraryPath_ErrorByService(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	expectedErrorMessage := "expected error message"
	s.mockService.LibraryService.MockErrors[0] = errors.New(expectedErrorMessage)
	r.POST("/", s.server.CreateLibrary)

	expectedName := "expectedLibraryName"
	req, err := http.NewRequest("POST", "/", body(fmt.Sprintf(`{"name":"%v"}`, expectedName)))
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

func Test_CreateLibraryPath_Success(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	expectedId, _ := uuid.NewRandom()
	expectedLibraryName := "some expected library name"
	s.mockService.LibraryService.MockModel[0] = &model.Library{
		ID:   expectedId,
		Name: expectedLibraryName,
	}

	r.POST("/", s.server.CreateLibrary)

	req, err := http.NewRequest("POST", "/", body(fmt.Sprintf(`{"name":"%v"}`, expectedLibraryName)))
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
