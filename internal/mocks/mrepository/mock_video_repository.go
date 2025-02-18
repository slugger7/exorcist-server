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

func (m *MockVideoRepo) UpdateExists(video *model.Video) error {
	stack := incStack()
	return m.MockError[stack]
}

func (m *MockVideoRepo) Insert(models []model.Video) ([]model.Video, error) {
	stack := incStack()
	return m.MockModels[stack], m.MockError[stack]
}

func (m *MockVideoRepo) GetById(id uuid.UUID) (*model.Video, error) {
	stack := incStack()
	return m.MockModel[stack], m.MockError[stack]
}

func (m *MockVideoRepo) GetByIdWithLibraryPath(id uuid.UUID) (*videoRepository.VideoLibraryPathModel, error) {
	panic("unimplemented")
}

// UpdateChecksum implements videoRepository.IVideoRepository.
func (m *MockVideoRepo) UpdateChecksum(video *model.Video) error {
	panic("unimplemented")
}
