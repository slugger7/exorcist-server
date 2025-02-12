package mrepository

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
)

type MockLibraryPathRepo struct {
	MockModel  map[int]*model.LibraryPath
	MockModels map[int][]model.LibraryPath
	MockError  map[int]error
}

func SetupMockLibraryPathRepository() *MockLibraryPathRepo {
	mockModels := make(map[int][]model.LibraryPath)
	mockModel := make(map[int]*model.LibraryPath)
	mockError := make(map[int]error)
	return &MockLibraryPathRepo{MockModel: mockModel, MockModels: mockModels, MockError: mockError}
}

func (mr MockRepo) LibraryPathRepo() libraryPathRepository.ILibraryPathRepository {
	return mr.MockLibraryPathRepo
}

func (mlp MockLibraryPathRepo) Create(string, uuid.UUID) (*model.LibraryPath, error) {
	stack := incStack()
	return mlp.MockModel[stack], mlp.MockError[stack]
}
func (mlp MockLibraryPathRepo) GetLibraryPaths() ([]model.LibraryPath, error) {
	stack := incStack()
	return mlp.MockModels[stack], mlp.MockError[stack]
}
