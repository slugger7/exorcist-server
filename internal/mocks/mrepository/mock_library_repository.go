package mrepository

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
)

type MockLibraryRepo mocks.MockFixture[model.Library]

func SetupMockLibraryRepo() *MockLibraryRepo {
	x := MockLibraryRepo(*mocks.SetupMockFixture[model.Library]())
	return &x
}

func (mr MockRepository) LibraryRepo() libraryRepository.ILibraryRepository {
	return mr.MockLibraryRepo
}

func (mlr MockLibraryRepo) CreateLibrary(name string) (*model.Library, error) {
	stack := incStack()
	return mlr.MockModel[stack], mlr.MockError[stack]
}
func (mlr MockLibraryRepo) GetLibraryByName(name string) (*model.Library, error) {
	stack := incStack()
	return mlr.MockModel[stack], mlr.MockError[stack]
}
func (mlr MockLibraryRepo) GetLibraries() ([]model.Library, error) {
	stack := incStack()
	return mlr.MockModels[stack], mlr.MockError[stack]
}
func (mlr MockLibraryRepo) GetLibraryById(uuid.UUID) (*model.Library, error) {
	stack := incStack()
	return mlr.MockModel[stack], mlr.MockError[stack]
}
