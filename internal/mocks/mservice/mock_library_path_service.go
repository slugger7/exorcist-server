package mservice

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	libraryPathService "github.com/slugger7/exorcist/internal/service/library_path"
)

// Deprecated: moved to mockgen in mock folder
type MockLibaryPathService mocks.MockFixture[model.LibraryPath]

// Deprecated: moved to mockgen in mock folder
func SetupMockLibraryPathService() *MockLibaryPathService {
	x := MockLibaryPathService(*mocks.SetupMockFixture[model.LibraryPath]())
	return &x
}

// Deprecated: moved to mockgen in mock folder
func (ms MockService) LibraryPath() libraryPathService.ILibraryPathService {
	return ms.libraryPath
}

// Deprecated: moved to mockgen in mock folder
func (lps MockLibaryPathService) Create(*model.LibraryPath) (*model.LibraryPath, error) {
	stack := incStack()
	return lps.MockModel[stack], lps.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (lps MockLibaryPathService) GetAll() ([]model.LibraryPath, error) {
	stack := incStack()
	return lps.MockModels[stack], lps.MockError[stack]
}
