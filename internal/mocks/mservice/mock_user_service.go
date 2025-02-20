package mservice

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	"github.com/slugger7/exorcist/internal/models"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

// Deprecated: moved to mockgen in mock folder
type MockUserService mocks.MockFixture[model.User]

// Deprecated: moved to mockgen in mock folder
func SetupMockUserService() *MockUserService {
	x := MockUserService(*mocks.SetupMockFixture[model.User]())
	return &x
}

// Deprecated: moved to mockgen in mock folder
func (ms MockService) User() userService.IUserService {
	return ms.user
}

// Deprecated: moved to mockgen in mock folder
func (mus MockUserService) Create(username, password string) (*model.User, error) {
	stack := incStack()
	return mus.MockModel[stack], mus.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (mus MockUserService) Validate(username, password string) (*model.User, error) {
	stack := incStack()
	return mus.MockModel[stack], mus.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (us MockUserService) UpdatePassword(id uuid.UUID, model models.ResetPasswordModel) error {
	panic("unimplemented")
}
