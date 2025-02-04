package userService

import (
	"errors"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

type mockRepo struct {
	mockUserRepo
}

var count = 0

type mockUserRepo struct {
	mockModels map[int]*model.User
	mockError  error
}

func (mr mockRepo) UserRepo() userRepository.IUserRepository {
	return mr.mockUserRepo
}

func (mur mockUserRepo) CreateUser(user model.User) (*model.User, error) {
	if len(mur.mockModels) > count {
		return mur.mockModels[count], mur.mockError
	}
	count = count + 1
	return nil, mur.mockError
}

func (mur mockUserRepo) GetUserByUsernameAndPassword(username, password string) (*model.User, error) {
	if len(mur.mockModels) > count {
		return mur.mockModels[count], mur.mockError
	}
	count = count + 1
	return nil, mur.mockError
}
func (mur mockUserRepo) GetUserByUsername(username string, columns ...postgres.Projection) (*model.User, error) {
	if len(mur.mockModels) > count {
		return mur.mockModels[count], mur.mockError
	}
	return nil, mur.mockError
}

func Test_UserExists_ErrorFromRepo(t *testing.T) {
	expectedErr := errors.New("expected error")
	us := &UserService{repo: mockRepo{mockUserRepo{mockError: expectedErr}}}

	if _, err := us.UserExists(""); err.Error() != expectedErr.Error() {
		t.Errorf("Encountered an unexpected error checking to see if a user exsists\nExpected: %v\nGot: %v", expectedErr.Error(), err.Error())
	}
}

func Test_UserExists_UserIsNil_ShouldReturnFalse(t *testing.T) {
	us := &UserService{repo: mockRepo{mockUserRepo{}}}

	actual, err := us.UserExists("")
	if err != nil {
		t.Fatal(err)
	}

	if actual {
		t.Error("Expected user to not exist")
	}
}

func Test_UserExists_UserIsDefined_ShouldReturnTrue(t *testing.T) {
	count = 0
	models := make(map[int]*model.User)
	models[0] = &model.User{}
	us := &UserService{repo: mockRepo{mockUserRepo{mockModels: models}}}

	actual, err := us.UserExists("")
	if err != nil {
		t.Fatal(err)
	}

	if !actual {
		t.Error("Expected user to exist")
	}
}

func (mr mockRepo) Health() map[string]string {
	panic("not implemented")
}
func (mr mockRepo) Close() error {
	panic("not implemented")
}
func (mr mockRepo) JobRepo() jobRepository.IJobRepository {
	panic("not implemented")
}
func (mr mockRepo) LibraryPathRepo() libraryPathRepository.ILibraryPathRepository {
	panic("not implemented")
}
func (mr mockRepo) VideoRepo() videoRepository.IVideoRepository {
	panic("not implemented")
}
func (mr mockRepo) LibraryRepo() libraryRepository.ILibraryRepository {
	panic("not implemented")
}
