package server

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func Test_CreateLibraryPath_InvalidBody(t *testing.T) {
	s := setupOldServer()

	rr := s.withPostEndpoint(s.server.CreateLibraryPath).
		withPostRequest(body(`{invalid json}`)).
		exec()

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
	s := setupOldServer()

	rr := s.withPostEndpoint(s.server.CreateLibraryPath).
		withPostRequest(body(`{"path": ""}`)).
		exec()

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
	s := setupOldServer()

	expectedId, _ := uuid.NewRandom()
	expectedLibraryId, _ := uuid.NewRandom()
	expectedLibraryPath := "some/expected/path"
	s.mockService.LibraryPath.MockModel[0] = &model.LibraryPath{
		ID:        expectedId,
		LibraryID: expectedLibraryId,
		Path:      expectedLibraryPath,
	}

	rr := s.withPostEndpoint(s.server.CreateLibraryPath).
		withPostRequest(body(`{"path":"%v", "libraryId": "%v"}`, expectedLibraryPath, expectedLibraryId)).
		exec()

	expectedStatusCode := http.StatusCreated
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nexpected %v but got %v", expectedStatusCode, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"ID":"%v","LibraryID":"%v","Path":"%v","Created":"0001-01-01T00:00:00Z","Modified":"0001-01-01T00:00:00Z"}`, expectedId.String(), expectedLibraryId.String(), expectedLibraryPath)
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("incorrect body\nexpected %v but got %v", expectedBody, body)
	}
}

func Test_GetAllLibraryPaths_WithServiceThrowingError(t *testing.T) {
	s := setupOldServer()

	expectedError := "expected error"
	s.mockService.LibraryPath.MockError[0] = errors.New(expectedError)

	rr := s.withGetEndpoint(s.server.GetAllLibraryPaths, "").
		withGetRequest(nil, "").
		exec()

	expectedStatusCode := http.StatusInternalServerError
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nExpected: %v Got: %v", expectedStatusCode, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"error":"%v"}`, ErrGetAllLibraryPathsService)
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("Expected: %v\nGot: %v", expectedBody, body)
	}
}

func Test_GetAllLibraryPaths_Success(t *testing.T) {
	s := setupOldServer()

	id, _ := uuid.NewRandom()
	libPath := model.LibraryPath{ID: id}
	libPaths := []model.LibraryPath{libPath}
	s.mockService.LibraryPath.MockModels[0] = libPaths

	rr := s.withGetEndpoint(s.server.GetAllLibraryPaths, "").
		withGetRequest(nil, "").
		exec()

	expectedStatusCode := http.StatusOK
	if rr.Code != expectedStatusCode {
		t.Errorf("wrong status code returned\nExpected: %v Got: %v", expectedStatusCode, rr.Code)
	}
	expectedBody := fmt.Sprintf(`[{"ID":"%v","LibraryID":"00000000-0000-0000-0000-000000000000","Path":"","Created":"0001-01-01T00:00:00Z","Modified":"0001-01-01T00:00:00Z"}]`, id)
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("Expected: %v\nGot: %v", expectedBody, body)
	}
}
