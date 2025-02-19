package mrepository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
)

// Deprecated: moved to mockgen in mock folder
type MockUserRepo mocks.MockFixture[model.User]

// Deprecated: moved to mockgen in mock folder
func SetupMockUserRepository() *MockUserRepo {
	x := MockUserRepo(*mocks.SetupMockFixture[model.User]())
	return &x
}

// Deprecated: moved to mockgen in mock folder
func (mr MockRepository) User() userRepository.IUserRepository {
	return mr.MockUserRepo
}

// Deprecated: moved to mockgen in mock folder
func (mur *MockUserRepo) GetUserByUsernameAndPassword(string, string) (*model.User, error) {
	stack := incStack()
	return mur.MockModel[stack], mur.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (mur *MockUserRepo) GetUserByUsername(string, ...postgres.Projection) (*model.User, error) {
	stack := incStack()
	return mur.MockModel[stack], mur.MockError[stack]
}

// Deprecated: moved to mockgen in mock folder
func (mur *MockUserRepo) CreateUser(model.User) (*model.User, error) {
	stack := incStack()
	return mur.MockModel[stack], mur.MockError[stack]
}
