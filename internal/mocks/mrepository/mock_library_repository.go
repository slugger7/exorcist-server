package mrepository

import "github.com/slugger7/exorcist/internal/db/exorcist/public/model"

type MockLibraryRepo struct {
	MockModel  map[int]*model.Library
	MockError  map[int]error
	MockModels map[int][]model.Library
}

func SetupMockLibraryRepository() *MockLibraryRepo {
	mockModels := make(map[int][]model.Library)
	mockModel := make(map[int]*model.Library)
	mockError := make(map[int]error)
	return &MockLibraryRepo{MockModel: mockModel, MockModels: mockModels, MockError: mockError}
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
