package mediaService

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	mock_repository "github.com/slugger7/exorcist/internal/mock/repository"
	mock_mediaRepository "github.com/slugger7/exorcist/internal/mock/repository/media"
	"github.com/slugger7/exorcist/internal/models"
	mediaRepository "github.com/slugger7/exorcist/internal/repository/media"
	"go.uber.org/mock/gomock"
)

type testService struct {
	svc            *mediaService
	repo           *mock_repository.MockRepository
	mediaRepo      *mock_mediaRepository.MockMediaRepository
	base           string
	mediaFile      string
	assetsPath     string
	assetId        uuid.UUID
	assetPath      string
	thumbnailAsset string
}

func setup(t *testing.T) *testService {
	base := "test_data"
	mediaFile := path.Join(base, "mediafile.mp4")
	assetsPath := path.Join(base, "assets")
	assetId, _ := uuid.NewRandom()
	assetPath := path.Join(assetsPath, assetId.String())
	thumbnailAsset := path.Join(assetsPath, "thumbnail.webp")
	ctrl := gomock.NewController(t)

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockMediaRepo := mock_mediaRepository.NewMockMediaRepository(ctrl)

	mockRepo.EXPECT().
		Media().DoAndReturn(func() mediaRepository.MediaRepository {
		return mockMediaRepo
	}).AnyTimes()

	ms := &mediaService{repo: mockRepo}

	if err := os.MkdirAll(assetsPath, os.ModePerm); err != nil {
		t.Fatalf("colud not create test data folder: %v", err.Error())
	}

	if _, err := os.Create(mediaFile); err != nil {
		t.Fatalf("could not create media file: %v", err.Error())
	}

	if _, err := os.Create(thumbnailAsset); err != nil {
		t.Fatalf("could not create thumbnail asset: %v", err.Error())
	}

	return &testService{ms, mockRepo, mockMediaRepo, base, mediaFile, assetsPath, assetId, assetPath, thumbnailAsset}
}

func (t *testService) cleanup() {
	os.RemoveAll(t.base)
}

func Test_Delete_PhysicalIsFalse_UpdatesDeletedFlagOnly(t *testing.T) {
	s := setup(t)
	defer s.cleanup()

	mediaEntity := models.Media{
		Media: model.Media{
			ID:      s.assetId,
			Exists:  true,
			Deleted: false,
			Path:    s.mediaFile,
		},
	}

	assets := []model.Media{
		{
			Exists:  true,
			Deleted: false,
		},
	}

	s.mediaRepo.EXPECT().
		GetById(s.assetId).
		DoAndReturn(func(i uuid.UUID) (*models.Media, error) {
			return &mediaEntity, nil
		}).Times(1)

	s.mediaRepo.EXPECT().
		GetAssetsFor(s.assetId).
		DoAndReturn(func(id uuid.UUID) ([]model.Media, error) {
			return assets, nil
		}).Times(1)

	s.mediaRepo.EXPECT().
		Delete(gomock.Any()).
		DoAndReturn(func(m model.Media) error {
			if m.Deleted != true {
				t.Errorf("Deleted was not set to true")
			}

			if m.Exists != true {
				t.Errorf("Exists was set to false when it was not a physical delete")
			}

			return nil
		}).Times(2)

	if err := s.svc.Delete(s.assetId, false); err != nil {
		t.Errorf("was not expecting an error but received: %v", err.Error())
	}
}

func Test_Delete_PhysicalIsTrue_UpdatesDeletedAndExistFlag_RemovesMediaAndAssets(t *testing.T) {
	s := setup(t)
	defer s.cleanup()

	mediaEntity := models.Media{
		Media: model.Media{
			ID:      s.assetId,
			Exists:  true,
			Deleted: false,
		},
	}
	assets := []model.Media{
		{
			Exists:  true,
			Deleted: false,
		},
	}

	s.mediaRepo.EXPECT().
		GetById(s.assetId).
		DoAndReturn(func(i uuid.UUID) (*models.Media, error) {
			return &mediaEntity, nil
		}).Times(1)

	s.mediaRepo.EXPECT().
		GetAssetsFor(s.assetId).
		DoAndReturn(func(id uuid.UUID) ([]model.Media, error) {
			return assets, nil
		}).Times(1)

	s.mediaRepo.EXPECT().
		Delete(gomock.Any()).
		DoAndReturn(func(m model.Media) error {
			if m.Deleted != true {
				t.Errorf("Deleted was not set to true")
			}

			if m.Exists == false {
				t.Errorf("Exists was was not set to false")
			}

			return nil
		}).Times(2)

	if err := s.svc.Delete(s.assetId, false); err != nil {
		t.Errorf("was not expecting an error but received: %v", err.Error())
	}

	if _, err := os.Stat(s.assetPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected error file not exist but got: %v", err.Error())
		}
	} else {
		t.Errorf("expected not exist error but file exists")
	}
}
