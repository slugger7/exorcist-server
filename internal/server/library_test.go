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
	s := &Server{}

	r.POST("/", s.CreateLibrary)

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

func Test_CreateLibrary_ErrorByService(t *testing.T) {
	r := setupEngine()
	s := &Server{}

	expectedErrorMessage := "expected error message"
	s.service = mockService{mockUserService{}, mockLibraryService{returningModel: nil, returningError: errors.New(expectedErrorMessage)}}
	r.POST("/", s.CreateLibrary)

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

func Test_CreateLibrary_Success(t *testing.T) {
	r := setupEngine()
	s := &Server{}

	expectedId, _ := uuid.NewRandom()

	expectedLibraryName := "some expected library name"
	s.service = mockService{mockUserService{}, mockLibraryService{returningModel: &model.Library{
		ID:   expectedId,
		Name: expectedLibraryName,
	}, returningError: nil}}
	r.POST("/", s.CreateLibrary)

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
	expectedBody := fmt.Sprintf(`{"ID":"%v","Name":"%v","Created":"0001-01-01T00:00:00Z","Modified":"0001-01-01T00:00:00Z"}`, expectedId.String(), expectedLibraryName)
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}
