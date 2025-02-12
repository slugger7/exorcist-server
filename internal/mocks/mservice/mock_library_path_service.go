package mservice

import "github.com/slugger7/exorcist/internal/db/exorcist/public/model"

type MockLibaryPathService struct {
	MockModels map[int][]model.LibraryPath
	MockModel  map[int]*model.LibraryPath
	MockError  map[int]error
}

func SetupMockLibraryPathService() MockLibaryPathService {
	mockModels := make(map[int][]model.LibraryPath)
	mockErrors := make(map[int]error)
	mockModel := make(map[int]*model.LibraryPath)
	return MockLibaryPathService{mockModels, mockModel, mockErrors}
}

func (lps MockLibaryPathService) Create(*model.LibraryPath) (*model.LibraryPath, error) {
	stack := incStack()
	return lps.MockModel[stack], lps.MockError[stack]
}

func (lps MockLibaryPathService) GetAll() ([]model.LibraryPath, error) {
	stack := incStack()
	return lps.MockModels[stack], lps.MockError[stack]
}
