package videoService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	mock_repository "github.com/slugger7/exorcist/internal/mock/repository"
	mock_videoRepository "github.com/slugger7/exorcist/internal/mock/repository/video"
	"github.com/slugger7/exorcist/internal/models"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
	"go.uber.org/mock/gomock"
)

type testService struct {
	svc       *VideoService
	repo      *mock_repository.MockIRepository
	videoRepo *mock_videoRepository.MockIVideoRepository
}

func setup(t *testing.T) *testService {
	ctrl := gomock.NewController(t)

	mockRepo := mock_repository.NewMockIRepository(ctrl)
	mockVideoRepo := mock_videoRepository.NewMockIVideoRepository(ctrl)

	mockRepo.EXPECT().
		Video().
		DoAndReturn(func() videoRepository.IVideoRepository {
			return mockVideoRepo
		}).
		AnyTimes()

	vs := &VideoService{repo: mockRepo}
	return &testService{vs, mockRepo, mockVideoRepo}
}

func Test_GetAll_ErrorFromRepo(t *testing.T) {
	s := setup(t)

	expectedErr := errors.New("error")

	s.videoRepo.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.Video, error) {
			return nil, expectedErr
		}).
		Times(1)

	vids, err := s.svc.GetAll()
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
	s := setup(t)

	id, _ := uuid.NewRandom()
	videos := []model.Video{{ID: id}}

	s.videoRepo.EXPECT().
		GetAll().
		DoAndReturn(func() ([]model.Video, error) {
			return videos, nil
		}).
		Times(1)

	vids, err := s.svc.GetAll()
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
	s := setup(t)

	id, _ := uuid.NewRandom()

	s.videoRepo.EXPECT().
		GetById(id).
		DoAndReturn(func(uuid.UUID) (*model.Video, error) {
			return nil, fmt.Errorf("err")
		}).
		Times(1)

	vid, err := s.svc.GetById(id)
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
	s := setup(t)

	id, _ := uuid.NewRandom()
	video := model.Video{ID: id}

	s.videoRepo.EXPECT().
		GetById(id).
		DoAndReturn(func(uuid.UUID) (*model.Video, error) {
			return &video, nil
		}).
		Times(1)

	vid, err := s.svc.GetById(id)
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

func Test_GetOverview_ErrorFromRepo(t *testing.T) {
	s := setup(t)

	s.videoRepo.EXPECT().
		GetOverview().
		DoAndReturn(func() ([]models.VideoOverviewModel, error) {
			return nil, fmt.Errorf("expected error")
		})

	vids, err := s.svc.GetOverview()
	if err != nil {
		var e errs.IError
		if errors.As(err, &e) {
			if e.Message() != ErrVideoRepoOverview {
				t.Errorf("expected: %v\ngot: %v", ErrVideoRepoOverview, e.Message())
			}
		}
	} else {
		t.Error("expected error but was nil")
	}

	if vids != nil {
		t.Fatalf("expected error but videos were defined %v", vids)
	}
}
