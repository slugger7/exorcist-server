package userService

import (
	"errors"
	"fmt"
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
	mockErrors map[int]error
}

func (mr mockRepo) UserRepo() userRepository.IUserRepository {
	return mr.mockUserRepo
}

func (mur mockUserRepo) CreateUser(user model.User) (*model.User, error) {
	count = count + 1
	return mur.mockModels[count-1], mur.mockErrors[count-1]
}

func (mur mockUserRepo) GetUserByUsernameAndPassword(username, password string) (*model.User, error) {
	count = count + 1
	return mur.mockModels[count-1], mur.mockErrors[count-1]
}
func (mur mockUserRepo) GetUserByUsername(username string, columns ...postgres.Projection) (*model.User, error) {
	count = count + 1
	return mur.mockModels[count-1], mur.mockErrors[count-1]
}

func beforeEach() (map[int]*model.User, map[int]error) {
	count = 0
	models := make(map[int]*model.User)
	errs := make(map[int]error)
	return models, errs
}

func Test_UserExists_ErrorFromRepo(t *testing.T) {
	_, errs := beforeEach()
	errs[0] = errors.New("expected error")
	us := &UserService{repo: mockRepo{mockUserRepo{mockErrors: errs}}}

	if _, err := us.UserExists(""); err.Error() != errs[0].Error() {
		t.Errorf("Encountered an unexpected error checking to see if a user exsists\nExpected: %v\nGot: %v", errs[0].Error(), err.Error())
	}
}

func Test_UserExists_UserIsNil_ShouldReturnFalse(t *testing.T) {
	models, errs := beforeEach()
	models[0] = nil
	errs[0] = nil
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
	models, _ := beforeEach()
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

func Test_CreateUser_UserExistsRaisesError_ShouldReturnError(t *testing.T) {
	_, errs := beforeEach()
	expectedError := errors.New("expected error")
	errs[0] = expectedError
	us := &UserService{repo: mockRepo{mockUserRepo{mockErrors: errs}}}
	username := ""
	wrapedExpectedError := errors.Join(errors.New(fmt.Sprintf("could not determine if user '%v' exists", username)), expectedError)
	if _, err := us.CreateUser(username, ""); err.Error() != wrapedExpectedError.Error() {
		t.Errorf("Unexpected error thrown\nExpected: %v\nGot: %v", wrapedExpectedError.Error(), err.Error())
	}
}

func Test_CreateUser_UserExistsTrue_ShouldReturnError(t *testing.T) {
	models, _ := beforeEach()
	expectedError := errors.New("user already exists")
	models[0] = &model.User{}
	us := &UserService{repo: mockRepo{mockUserRepo{mockModels: models}}}
	username := ""

	if _, err := us.CreateUser(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Unexpected error thrown\nExpected: %v\nGot: %v", expectedError.Error(), err.Error())
	}
}

func Test_CreateUser_UserExistsFalse_RepoCreateReturnsError_ShouldReturnError(t *testing.T) {
	models, errs := beforeEach()
	expectedError := errors.New("expected error")
	errs[1] = expectedError
	username := ""
	wrapedExpectedError := errors.Join(errors.New("could not create a new user"), expectedError)
	us := &UserService{repo: mockRepo{mockUserRepo{mockModels: models, mockErrors: errs}}}

	if _, err := us.CreateUser(username, ""); err.Error() != wrapedExpectedError.Error() {
		t.Errorf("Unexpected error thrown\nExpected: %v\nGot: %v", wrapedExpectedError.Error(), err.Error())
	}
}

func Test_CreateUser_UserExistsFalse_RepoCreatesUser_ShouldReturnUser(t *testing.T) {
	models, errs := beforeEach()
	username := "someUser"
	password := "somePassword"
	models[1] = &model.User{Username: username, Password: password}
	us := &UserService{repo: mockRepo{mockUserRepo{mockModels: models, mockErrors: errs}}}

	user, err := us.CreateUser(username, "")
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Error("Created user is nil")
	}
	if user.Username != username {
		t.Errorf("Unexpected user returned\nExpected: %v\nGot:%v", username, user.Username)
	}
	if user.Password != "" {
		t.Error("Password was not removed before returning")
	}
}

// Unused mocks

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
