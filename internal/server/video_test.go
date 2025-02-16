package server

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func Test_GetVideo_InvalidId(t *testing.T) {
	s := setupServer()

	invalidId := "someinvalidid"

	rr := s.withGetEndpoint(s.server.GetVideo, ":id").
		withGetRequest(nil, invalidId).
		exec()

	expectedStatus := http.StatusBadRequest
	if expectedStatus != rr.Code {
		t.Errorf("Expected status code: %v\nGot status code: %v", expectedStatus, rr.Code)
	}
	expectedBody := fmt.Sprintf(`{"error":"%v"}`, ErrInvalidIdFormat)
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body: %v\nGot body: %v", expectedBody, rr.Body.String())
	}
}

func Test_GetVideo_ServiceReturnsError(t *testing.T) {
	s := setupServer()

	id, _ := uuid.NewRandom()

	s.mockService.Video.MockError[0] = fmt.Errorf("error")

	rr := s.withGetEndpoint(s.server.GetVideo, ":id").
		withGetRequest(nil, id.String()).
		exec()

	expectedStatus := http.StatusInternalServerError
	if expectedStatus != rr.Code {
		t.Errorf("Expected status code: %v\nGot status code: %v", expectedStatus, rr.Code)
	}

	expectedBody := fmt.Sprintf(`{"error":"%v"}`, ErrGetVideoService)
	if expectedBody != rr.Body.String() {
		t.Errorf("Expected body: %v\nGot body: %v", expectedBody, rr.Body.String())
	}
}

func Test_GetVideo_VideoServiceNil(t *testing.T) {
	s := setupServer()

	id, _ := uuid.NewRandom()

	s.mockService.Video.MockModel[0] = nil

	rr := s.withGetEndpoint(s.server.GetVideo, ":id").
		withGetRequest(nil, id.String()).
		exec()

	expectedStatus := http.StatusNotFound
	if expectedStatus != rr.Code {
		t.Errorf("Expected status code: %v\nGot status code: %v", expectedStatus, rr.Code)
	}

	expectedBody := fmt.Sprintf(`{"error":"%v"}`, ErrVideoNotFound)
	if expectedBody != rr.Body.String() {
		t.Errorf("Expected body: %v\nGot body: %v", expectedBody, rr.Body.String())
	}
}

func Test_GetVideo_Success(t *testing.T) {
	s := setupServer()

	id, _ := uuid.NewRandom()
	video := &model.Video{ID: id}
	s.mockService.Video.MockModel[0] = video

	rr := s.withGetEndpoint(s.server.GetVideo, ":id").
		withGetRequest(nil, id.String()).
		exec()

	expectedStatus := http.StatusOK
	if expectedStatus != rr.Code {
		t.Errorf("Expected status code: %v\nGot status code: %v", expectedStatus, rr.Code)
	}

	expectedBody := fmt.Sprintf(`{"ID":"%v","LibraryPathID":"00000000-0000-0000-0000-000000000000","RelativePath":"","Title":"","FileName":"","Height":0,"Width":0,"Runtime":0,"Size":0,"Checksum":null,"Added":"0001-01-01T00:00:00Z","Deleted":false,"Exists":false,"Created":"0001-01-01T00:00:00Z","Modified":"0001-01-01T00:00:00Z"}`, id)
	if expectedBody != rr.Body.String() {
		t.Errorf("Expected body: %v\nGot body: %v", expectedBody, rr.Body.String())
	}
}
