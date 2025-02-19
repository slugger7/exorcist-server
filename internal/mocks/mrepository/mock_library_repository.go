package mrepository

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
)

// Deprecated: moved to mockgen in mock folder
type MockLibraryRepo mocks.MockFixture[model.Library]

// Deprecated: moved to mockgen in mock folder
func SetupMockLibraryRepo() *MockLibraryRepo {
	x := MockLibraryRepo(*mocks.SetupMockFixture[model.Library]())
	return &x
}

// Deprecated: moved to mockgen in mock folder
func (mr MockRepository) Library() libraryRepository.ILibraryRepository {
	return mr.MockLibraryRepo
}

// Deprecated: moved to mockgen in mock folder
func (mlr MockLibraryRepo) CreateLibrary(name string) (*model.Library, error) {
	stack := incStack()
	return mlr.MockModel[stack], mlr.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (mlr MockLibraryRepo) GetLibraryByName(name string) (*model.Library, error) {
	stack := incStack()
	return mlr.MockModel[stack], mlr.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (mlr MockLibraryRepo) GetLibraries() ([]model.Library, error) {
	stack := incStack()
	return mlr.MockModels[stack], mlr.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (mlr MockLibraryRepo) GetLibraryById(uuid.UUID) (*model.Library, error) {
	stack := incStack()
	return mlr.MockModel[stack], mlr.MockError[stack]
}
