package mservice

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
)

type MockLibraryService struct {
	MockModels map[int][]model.Library
	MockErrors map[int]error
	MockModel  map[int]*model.Library
}

func (ls MockLibraryService) CreateLibrary(actual model.Library) (*model.Library, error) {
	stack := incStack()
	return ls.MockModel[stack], ls.MockErrors[stack]
}

func (ls MockLibraryService) GetLibraries() ([]model.Library, error) {
	stack := incStack()
	return ls.MockModels[stack], ls.MockErrors[stack]
}

func SetupMockLibraryService() MockLibraryService {
	mockModels := make(map[int][]model.Library)
	mockErrors := make(map[int]error)
	mockModel := make(map[int]*model.Library)
	return MockLibraryService{MockModels: mockModels, MockErrors: mockErrors, MockModel: mockModel}
}
