package server

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

func Test_GetVideo_InvalidIdD(t *testing.T) {
	s := setupServer(t)

	id := "some invalid uuid"

	rr := s.withGetEndpoint(s.server.GetVideo, ":id").
		withGetRequest(id).
		exec()

	assert.StatusCode(t, http.StatusBadRequest, rr.Code)
	assert.Body(t, fmt.Sprintf(`{"error":"%v"}`, ErrInvalidIdFormat), rr.Body.String())
}

func Test_GetVideo_ServiceReturnsError(t *testing.T) {
	s := setupOldServer()

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
	s := setupOldServer()

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
	s := setupOldServer()

	id, _ := uuid.NewRandom()
	video := &model.Video{ID: id}
	s.mockService.Video.MockModel[0] = video

	rr := s.withGetEndpoint(s.server.GetVideo, ":id").
		withGetRequest(nil, id.String()).
		exec()

	expectedStatus := http.StatusMovedPermanently
	if expectedStatus != rr.Code {
		t.Errorf("Expected status code: %v\nGot status code: %v", expectedStatus, rr.Code)
	}
	// TODO: figure out how to test that a video is returned
}
