package mrepository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/mocks"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
)

type MockUserRepo mocks.MockFixture[model.User]

func SetupMockUserRepository() *MockUserRepo {
	x := MockUserRepo(*mocks.SetupMockFixture[model.User]())
	return &x
}

func (mr MockRepository) User() userRepository.IUserRepository {
	return mr.MockUserRepo
}

func (mur *MockUserRepo) GetUserByUsernameAndPassword(string, string) (*model.User, error) {
	stack := incStack()
	return mur.MockModel[stack], mur.MockError[stack]
}

func (mur *MockUserRepo) GetUserByUsername(string, ...postgres.Projection) (*model.User, error) {
	stack := incStack()
	return mur.MockModel[stack], mur.MockError[stack]
}

func (mur *MockUserRepo) CreateUser(model.User) (*model.User, error) {
	stack := incStack()
	return mur.MockModel[stack], mur.MockError[stack]
}
