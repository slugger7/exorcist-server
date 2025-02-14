package mservice

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	userService "github.com/slugger7/exorcist/internal/service/user"
)

type MockUserService mocks.MockFixture[model.User]

func SetupMockUserService() *MockUserService {
	x := MockUserService(*mocks.SetupMockFixture[model.User]())
	return &x
}

func (ms MockService) User() userService.IUserService {
	return ms.user
}

func (mus MockUserService) Create(username, password string) (*model.User, error) {
	stack := incStack()
	return mus.MockModel[stack], mus.MockError[stack]
}

func (mus MockUserService) Validate(username, password string) (*model.User, error) {
	stack := incStack()
	return mus.MockModel[stack], mus.MockError[stack]
}
