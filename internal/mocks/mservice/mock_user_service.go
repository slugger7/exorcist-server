package mservice

import "github.com/slugger7/exorcist/internal/db/exorcist/public/model"

type MockUserService struct {
	MockModels map[int][]model.User
	MockErrors map[int]error
	MockModel  map[int]*model.User
}

func (mus MockUserService) CreateUser(username, password string) (*model.User, error) {
	stack := incStack()
	return mus.MockModel[stack], mus.MockErrors[stack]
}

func (mus MockUserService) ValidateUser(username, password string) (*model.User, error) {
	stack := incStack()
	return mus.MockModel[stack], mus.MockErrors[stack]
}

func SetupMockUserService() MockUserService {
	mockModels := make(map[int][]model.User)
	mockErrors := make(map[int]error)
	mockModel := make(map[int]*model.User)
	return MockUserService{MockModels: mockModels, MockErrors: mockErrors, MockModel: mockModel}
}
