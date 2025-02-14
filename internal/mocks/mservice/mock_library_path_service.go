package mservice

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	libraryPathService "github.com/slugger7/exorcist/internal/service/library_path"
)

type MockLibaryPathService mocks.MockFixture[model.LibraryPath]

func SetupMockLibraryPathService() *MockLibaryPathService {
	x := MockLibaryPathService(*mocks.SetupMockFixture[model.LibraryPath]())
	return &x
}

func (ms MockService) LibraryPath() libraryPathService.ILibraryPathService {
	return ms.libraryPath
}

func (lps MockLibaryPathService) Create(*model.LibraryPath) (*model.LibraryPath, error) {
	stack := incStack()
	return lps.MockModel[stack], lps.MockError[stack]
}

func (lps MockLibaryPathService) GetAll() ([]model.LibraryPath, error) {
	stack := incStack()
	return lps.MockModels[stack], lps.MockError[stack]
}
