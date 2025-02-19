package mservice

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
	videoService "github.com/slugger7/exorcist/internal/service/video"
)

type MockVideoService mocks.MockFixture[model.Video]

func SetupMockVideoService() *MockVideoService {
	x := MockVideoService(*mocks.SetupMockFixture[model.Video]())
	return &x
}

func (ms *MockService) Video() videoService.IVideoService {
	return ms.video
}

func (mvs *MockVideoService) GetAll() ([]model.Video, error) {
	stack := incStack()
	return mvs.MockModels[stack], mvs.MockError[stack]
}

func (mvs *MockVideoService) GetById(uuid.UUID) (*model.Video, error) {
	stack := incStack()
	return mvs.MockModel[stack], mvs.MockError[stack]
}

func (mvs *MockVideoService) GetByIdWithLibraryPath(uuid.UUID) (*videoRepository.VideoLibraryPathModel, error) {
	panic("todo")
}
