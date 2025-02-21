package mrepository

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

// Deprecated: moved to mockgen in mock folder
type MockVideoRepo mocks.MockFixture[model.Video]

// Deprecated: moved to mockgen in mock folder
func SetupMockVideoRepository() *MockVideoRepo {
	x := MockVideoRepo(*mocks.SetupMockFixture[model.Video]())
	return &x
}

// Deprecated: moved to mockgen in mock folder
func (mr *MockRepository) Video() videoRepository.IVideoRepository {
	return mr.MockVideoRepo
}

// Deprecated: moved to mockgen in mock folder
func (m *MockVideoRepo) GetAll() ([]model.Video, error) {
	stack := incStack()
	return m.MockModels[stack], m.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (m *MockVideoRepo) GetByLibraryPathId(id uuid.UUID) ([]model.Video, error) {
	stack := incStack()
	return m.MockModels[stack], m.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (m *MockVideoRepo) UpdateExists(video *model.Video) error {
	stack := incStack()
	return m.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (m *MockVideoRepo) Insert(models []model.Video) ([]model.Video, error) {
	stack := incStack()
	return m.MockModels[stack], m.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (m *MockVideoRepo) GetById(id uuid.UUID) (*model.Video, error) {
	stack := incStack()
	return m.MockModel[stack], m.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (m *MockVideoRepo) GetByIdWithLibraryPath(id uuid.UUID) (*videoRepository.VideoLibraryPathModel, error) {
	panic("unimplemented")
}

// Deprecated: moved to mockgen in mock folder
func (m *MockVideoRepo) UpdateChecksum(video *model.Video) error {
	panic("unimplemented")
}
