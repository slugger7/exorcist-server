package mrepository

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	imageRepository "github.com/slugger7/exorcist/internal/repository/image"
)

// Deprecated
type MockImageRepo mocks.MockFixture[model.Image]

// Deprecated
func (m *MockImageRepo) Create(model *model.Image) (*model.Image, error) {
	panic("unimplemented")
}

// Deprecated: use mockgen instead
func SetupMockImageRepo() *MockImageRepo {
	x := MockImageRepo(*mocks.SetupMockFixture[model.Image]())
	return &x
}

// Deprecated
func (mr MockRepository) Image() imageRepository.IImageRepository {
	return mr.MockImageRepo
}
