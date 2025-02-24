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

func Test_CreateLibraryPath_InvalidBody(t *testing.T) {
	s := setupServer(t)

	s.server.withLibraryPathCreate(&s.engine.RouterGroup, "/")
	rr := s.withPostRequest(body(`{some invalid body}`)).
		exec()

	assert.StatusCode(t, http.StatusUnprocessableEntity, rr.Code)
}

func Test_CreateLibraryPath_ErrFromService(t *testing.T) {
	s := setupServer(t).
		withLibraryPathService()

	id, _ := uuid.NewRandom()
	m := &models.CreateLibraryPathModel{
		LibraryId: id,
		Path:      "some path",
	}

	s.mockLibraryPathService.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(*model.LibraryPath) (*model.LibraryPath, error) {
			return nil, fmt.Errorf("some error")
		}).
		Times(1)

	s.server.withLibraryPathCreate(&s.engine.RouterGroup, "/")
	rr := s.withPostRequest(bodyM(m)).
		exec()

	assert.StatusCode(t, http.StatusInternalServerError, rr.Code)
	assert.Body(t, errBody(ErrCreatingLibraryPath), rr.Body.String())
}

func Test_CreateLibraryPath_Success(t *testing.T) {
	s := setupServer(t).
		withLibraryPathService()

	libId, _ := uuid.NewRandom()
	id, _ := uuid.NewRandom()
	m := &models.CreateLibraryPathModel{
		LibraryId: libId,
		Path:      "some path",
	}

	s.mockLibraryPathService.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(model *model.LibraryPath) (*model.LibraryPath, error) {
			model.ID = id
			return model, nil
		}).
		Times(1)

	s.server.withLibraryPathCreate(&s.engine.RouterGroup, "/")
	rr := s.withPostRequest(bodyM(m)).
		exec()

	result := models.LibraryPathModel{
		Id:        id,
		LibraryId: libId,
		Path:      m.Path,
	}

	body, _ := json.Marshal(result)

	assert.StatusCode(t, http.StatusCreated, rr.Code)
	assert.Body(t, string(body), rr.Body.String())
}

func Test_GetAllLibraryPaths_WithServiceThrowingError(t *testing.T) {
	s := setupServer(t).
		withLibraryPathService()

	s.mockLibraryPathService.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.LibraryPath, error) {
			return nil, fmt.Errorf("some error")
		})

	s.server.withLibraryPathGetAll(&s.engine.RouterGroup, "/")
	rr := s.withGetRequest("").
		exec()
	assert.StatusCode(t, http.StatusInternalServerError, rr.Code)
	assert.Body(t, errBody(ErrGetAllLibraryPathsService), rr.Body.String())
}

func Test_GetAllLibraryPaths_Success(t *testing.T) {
	s := setupServer(t).
		withLibraryPathService()

	id, _ := uuid.NewRandom()
	libId, _ := uuid.NewRandom()
	libPath := model.LibraryPath{
		ID:        id,
		LibraryID: libId,
		Path:      "some path",
	}

	s.mockLibraryPathService.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.LibraryPath, error) {
			return []model.LibraryPath{libPath}, nil
		})

	s.server.withLibraryPathGetAll(&s.engine.RouterGroup, "/")
	rr := s.withGetRequest("").
		exec()

	body, _ := json.Marshal([]models.LibraryPathModel{{Id: libPath.ID, LibraryId: libPath.LibraryID, Path: libPath.Path}})

	assert.StatusCode(t, http.StatusOK, rr.Code)
	assert.Body(t, string(body), rr.Body.String())

}
