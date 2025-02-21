package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/models"
	"go.uber.org/mock/gomock"
)

func Test_CreateLibrary_InvalidBody(t *testing.T) {
	s := setupServer(t)

	rr := s.withPostEndpoint(s.server.CreateLibrary).
		withPostRequest(body(`{invalid json}`)).
		exec()

	assert.StatusCode(t, http.StatusUnprocessableEntity, rr.Code)
}

func Test_CreateLibrary_ErrorByService(t *testing.T) {
	s := setupServer(t).
		withLibraryService()

	m := models.CreateLibraryModel{
		Name: "someName",
	}
	l := &model.Library{Name: m.Name}

	s.mockLibraryService.EXPECT().
		Create(gomock.Eq(l)).
		DoAndReturn(func(*model.Library) (*model.Library, error) {
			return nil, fmt.Errorf("some error")
		}).
		Times(1)

	rr := s.withPostEndpoint(s.server.CreateLibrary).
		withPostRequest(bodyM(m)).
		exec()

	assert.StatusCode(t, http.StatusBadRequest, rr.Code)
	assert.Body(t, errBody(ErrCreatingLibrary), rr.Body.String())
}

func Test_CreateLibrary_Success(t *testing.T) {
	s := setupServer(t).
		withLibraryService()

	m := models.CreateLibraryModel{
		Name: "someName",
	}
	l := model.Library{Name: m.Name}
	id, _ := uuid.NewRandom()

	s.mockLibraryService.EXPECT().
		Create(gomock.Eq(&l)).
		DoAndReturn(func(ml *model.Library) (*model.Library, error) {
			ml.ID = id
			return ml, nil
		}).
		Times(1)

	cm := model.Library{ID: id, Name: m.Name}

	rr := s.withPostEndpoint(s.server.CreateLibrary).
		withPostRequest(bodyM(m)).
		exec()

	result := (&models.Library{}).FromModel(cm)

	body, _ := json.Marshal(result)
	assert.StatusCode(t, http.StatusCreated, rr.Code)
	assert.Body(t, string(body), rr.Body.String())
}

func Test_GetLibraries_ServiceReturnsError(t *testing.T) {
	s := setupServer(t).
		withLibraryService()

	s.mockLibraryService.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.Library, error) {
			return nil, fmt.Errorf("some error")
		})

	rr := s.withGetEndpoint(s.server.GetLibraries, "").
		withGetRequest("").
		exec()

	assert.StatusCode(t, http.StatusInternalServerError, rr.Code)
	assert.Body(t, errBody(ErrGetLibraries), rr.Body.String())
}

func Test_GetLibraries_Succeeds(t *testing.T) {
	s := setupServer(t).
		withLibraryService()

	lib := model.Library{Name: "lib"}
	libs := []model.Library{lib}

	s.mockLibraryService.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.Library, error) {
			return libs, nil
		}).
		Times(1)

	rr := s.withGetEndpoint(s.server.GetLibraries, "").
		withGetRequest("").
		exec()

	bm := []models.Library{{Name: lib.Name}}

	body, _ := json.Marshal(bm)

	assert.StatusCode(t, http.StatusOK, rr.Code)
	assert.Body(t, string(body), rr.Body.String())
}

func Test_LibraryAction_WithInvalidId(t *testing.T) {
	s := setupOldServer()

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
	s := setupOldServer()

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
