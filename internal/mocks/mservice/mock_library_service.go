package mservice

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
)

type MockLibraryService mocks.MockFixture[model.Library]

func SetupMockLibraryService() *MockLibraryService {
	x := MockLibraryService(*mocks.SetupMockFixture[model.Library]())
	return &x
}

func (ms MockService) Library() libraryService.ILibraryService {
	return ms.library
}

func (ls MockLibraryService) Create(actual model.Library) (*model.Library, error) {
	stack := incStack()
	return ls.MockModel[stack], ls.MockError[stack]
}

func (ls MockLibraryService) GetAll() ([]model.Library, error) {
	stack := incStack()
	return ls.MockModels[stack], ls.MockError[stack]
}

func (ls MockLibraryService) Action(id uuid.UUID, action string) error {
	stack := incStack()
	return ls.MockError[stack]
}
