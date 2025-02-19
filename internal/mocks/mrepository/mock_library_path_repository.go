package mrepository

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
)

// Deprecated: moved to mockgen in mock folder
type MockLibraryPathRepo mocks.MockFixture[model.LibraryPath]

// Deprecated: moved to mockgen in mock folder
func SetupMockLibraryPathRepository() *MockLibraryPathRepo {
	x := MockLibraryPathRepo(*mocks.SetupMockFixture[model.LibraryPath]())
	return &x
}

// Deprecated: moved to mockgen in mock folder
func (mr MockRepository) LibraryPath() libraryPathRepository.ILibraryPathRepository {
	return mr.MockLibraryPathRepo
}

// Deprecated: moved to mockgen in mock folder
func (mlp MockLibraryPathRepo) Create(string, uuid.UUID) (*model.LibraryPath, error) {
	stack := incStack()
	return mlp.MockModel[stack], mlp.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (mlp MockLibraryPathRepo) GetAll() ([]model.LibraryPath, error) {
	stack := incStack()
	return mlp.MockModels[stack], mlp.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (mlp *MockLibraryPathRepo) GetByLibraryId(libraryId uuid.UUID) ([]model.LibraryPath, error) {
	stack := incStack()
	return mlp.MockModels[stack], mlp.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (mlp *MockLibraryPathRepo) GetById(id uuid.UUID) (*model.LibraryPath, error) {
	stack := incStack()
	return mlp.MockModel[stack], mlp.MockError[stack]
}
