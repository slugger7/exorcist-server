package mservice

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	libraryService "github.com/slugger7/exorcist/internal/service/library"
)

// Deprecated: moved to mockgen in mock folder
type MockLibraryService mocks.MockFixture[model.Library]

// Deprecated: moved to mockgen in mock folder
func SetupMockLibraryService() *MockLibraryService {
	x := MockLibraryService(*mocks.SetupMockFixture[model.Library]())
	return &x
}

// Deprecated: moved to mockgen in mock folder
func (ms MockService) Library() libraryService.ILibraryService {
	return ms.library
}

// Deprecated: moved to mockgen in mock folder
func (ls MockLibraryService) Create(actual *model.Library) (*model.Library, error) {
	stack := incStack()
	return ls.MockModel[stack], ls.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (ls MockLibraryService) GetAll() ([]model.Library, error) {
	stack := incStack()
	return ls.MockModels[stack], ls.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (ls MockLibraryService) Action(id uuid.UUID, action string) error {
	stack := incStack()
	return ls.MockError[stack]
}
