package mrepository

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

type MockVideoRepo mocks.MockFixture[model.Video]

func SetupMockVideoRepository() *MockVideoRepo {
	x := MockVideoRepo(*mocks.SetupMockFixture[model.Video]())
	return &x
}

func (mr *MockRepository) Video() videoRepository.IVideoRepository {
	return mr.MockVideoRepo
}

func (m *MockVideoRepo) GetAll() ([]model.Video, error) {
	stack := incStack()
	return m.MockModels[stack], m.MockError[stack]
}

func (m *MockVideoRepo) GetByLibraryPathId(id uuid.UUID) ([]model.Video, error) {
	stack := incStack()
	return m.MockModels[stack], m.MockError[stack]
}

func (m *MockVideoRepo) UpdateVideoExists(video model.Video) error {
	stack := incStack()
	return m.MockError[stack]
}

func (m *MockVideoRepo) Insert(models []model.Video) error {
	stack := incStack()
	return m.MockError[stack]
}
