package mservice

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
)

type MockLibraryService mocks.MockFixture[model.Library]

func (ls MockLibraryService) Create(actual model.Library) (*model.Library, error) {
	stack := incStack()
	return ls.MockModel[stack], ls.MockError[stack]
}

func (ls MockLibraryService) GetAll() ([]model.Library, error) {
	stack := incStack()
	return ls.MockModels[stack], ls.MockError[stack]
}

func SetupMockLibraryService() MockLibraryService {
	mockModels := make(map[int][]model.Library)
	mockErrors := make(map[int]error)
	mockModel := make(map[int]*model.Library)
	return MockLibraryService{MockModels: mockModels, MockError: mockErrors, MockModel: mockModel}
}
