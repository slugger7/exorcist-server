package mservice

import (
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
)

type MockUserService mocks.MockFixture[model.User]

func (mus MockUserService) Create(username, password string) (*model.User, error) {
	stack := incStack()
	return mus.MockModel[stack], mus.MockError[stack]
}

func (mus MockUserService) Validate(username, password string) (*model.User, error) {
	stack := incStack()
	return mus.MockModel[stack], mus.MockError[stack]
}

func SetupMockUserService() MockUserService {
	mockModels := make(map[int][]model.User)
	mockErrors := make(map[int]error)
	mockModel := make(map[int]*model.User)
	return MockUserService{MockModels: mockModels, MockError: mockErrors, MockModel: mockModel}
}
