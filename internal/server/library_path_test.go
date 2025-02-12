package server

import (
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

func Test_CreateLibraryPath_Success(t *testing.T) {
	r := setupEngine()
	s := setupServer()

	expectedId, _ := uuid.NewRandom()
	expectedLibraryId, _ := uuid.NewRandom()
	expectedLibraryPath := "some/expected/path"
	s.mockService.LibraryPathService.MockModel[0] = &model.LibraryPath{
		ID:        expectedId,
		LibraryID: expectedLibraryId,
		Path:      expectedLibraryPath,
	}

	r.POST("/", s.server.CreateLibraryPath)

	req, err := http.NewRequest("POST", "/", body(`{"path":"%v", "libraryId": "%v"}`, expectedLibraryPath, expectedLibraryId))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	expectedStatusCode := http.StatusCreated
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"ID":"%v","LibraryID":"%v","Path":"%v","Created":"0001-01-01T00:00:00Z","Modified":"0001-01-01T00:00:00Z"}`, expectedId.String(), expectedLibraryId.String(), expectedLibraryPath)
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}
