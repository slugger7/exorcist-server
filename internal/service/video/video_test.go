package videoService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	mock_repository "github.com/slugger7/exorcist/internal/mock/repository"
	mock_videoRepository "github.com/slugger7/exorcist/internal/mock/repository/video"
	"github.com/slugger7/exorcist/internal/mocks/mrepository"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
	"go.uber.org/mock/gomock"
)

func setup() (*VideoService, *mrepository.MockRepository) {
	mockRepo := mrepository.SetupMockRespository()
	vs := &VideoService{repo: mockRepo}
	return vs, mockRepo
}

func Test_GetAll_ErrorFromRepo(t *testing.T) {
	vs, mr := setup()

	mr.MockVideoRepo.MockError[0] = errors.New("error")

	vids, err := vs.GetAll()
	if err == nil {
		t.Error("Expected error but got nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		if e.Message() != ErrGetAllVideos {
			t.Errorf("Expected error: %v\nGot error: %v", ErrGetAllVideos, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got: %v", err)
	}

	if vids != nil {
		t.Fatalf("Vids was supposed to be nil but was: %v", vids)
	}
}

func Test_GetAll_Success(t *testing.T) {
	vs, mr := setup()

	id, _ := uuid.NewRandom()
	videos := []model.Video{{ID: id}}
	mr.MockVideoRepo.MockModels[0] = videos

	vids, err := vs.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(vids) != 1 {
		t.Error("Incorrect videos length")
	}
	actualId := vids[len(vids)-1].ID
	if actualId != id {
		t.Errorf("Expected video with id: %v\nGot video with id: %v", id, actualId)
	}
}

func Test_GetById_RepoReturnsError(t *testing.T) {
	vs, mr := setup()

	id, _ := uuid.NewRandom()

	mr.MockVideoRepo.MockError[0] = fmt.Errorf("err")

	vid, err := vs.GetById(id)
	if err == nil {
		t.Error("Expected error but got nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		expectedMessage := fmt.Sprintf(ErrVideoById, id)
		if e.Message() != expectedMessage {
			t.Errorf("Expected error: %v\nGot error: %v", expectedMessage, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got %v", err.Error())
	}

	if vid != nil {
		t.Errorf("Expected vid to be nil but was %v", vid)
	}
}

func Test_GetById_RepoReturnsVideo(t *testing.T) {
	vs, mr := setup()

	id, _ := uuid.NewRandom()
	video := model.Video{ID: id}
	mr.MockVideoRepo.MockModel[0] = &video

	vid, err := vs.GetById(id)
	if err != nil {
		t.Errorf("Expected nil but got %v", err.Error())
	}
	if vid == nil {
		t.Error("Expected video but was nil")
	}
	if vid.ID != id {
		t.Errorf("Expected video with id: %v\nGot video with id: %v", id, vid.ID)
	}
}

func Test_NewMocks(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := mock_videoRepository.NewMockIVideoRepository(ctrl)

	repo := mock_repository.NewMockIRepository(ctrl)
	repo.EXPECT().
		Video().
		DoAndReturn(func() IVideoService { return m }).
		AnyTimes()

	id, _ := uuid.NewRandom()

	m.EXPECT().
		GetByIdWithLibraryPath(gomock.Eq(id)).
		DoAndReturn(func(_ uuid.UUID) (*videoRepository.VideoLibraryPathModel, error) {
			return nil, nil
		}).
		AnyTimes()

	vs := New(repo, &environment.EnvironmentVariables{LogLevel: "none"})
	val, _ := vs.GetByIdWithLibraryPath(id)

	if val != nil {
		t.Error("This is an example method on how to use the mocks")
	}
}
