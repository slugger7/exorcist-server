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

func setupMockRepo() mockRepo {
	return mockRepo{mockUserRepo: mockUserRepo{
		mockModels: make(map[int]*model.User),
		mockErrors: make(map[int]error),
	}}
}

func setupMockUserService() (*UserService, *mockRepo) {
	count = 0
	mr := setupMockRepo()
	us := &UserService{repo: &mr}
	return us, &mr
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

func Test_UserExists_ErrorFromRepo(t *testing.T) {
	us, mr := setupMockUserService()

	expectedErr := errors.New("expected error")
	mr.mockUserRepo.mockErrors[0] = expectedErr

	if _, err := us.UserExists(""); err.Error() != expectedErr.Error() {
		t.Errorf("Expected: %v\nGot: %v", expectedErr.Error(), err.Error())
	}
}

func Test_UserExists_UserIsNil_ShouldReturnFalse(t *testing.T) {
	us, mr := setupMockUserService()
	mr.mockUserRepo.mockModels[0] = nil
	mr.mockUserRepo.mockErrors[0] = nil

	actual, err := us.UserExists("")
	if err != nil {
		t.Fatal(err)
	}

	if actual {
		t.Error("Expected user to not exist")
	}
}

func Test_UserExists_UserIsDefined_ShouldReturnTrue(t *testing.T) {
	us, mr := setupMockUserService()
	mr.mockUserRepo.mockModels[0] = &model.User{}

	actual, err := us.UserExists("")
	if err != nil {
		t.Fatal(err)
	}

	if !actual {
		t.Error("Expected user to exist")
	}
}

func Test_CreateUser_UserExistsRaisesError_ShouldReturnError(t *testing.T) {
	us, mr := setupMockUserService()
	expectedError := errors.New("expected error")
	mr.mockUserRepo.mockErrors[0] = expectedError
	username := ""
	expectedErrorMessage := fmt.Sprintf("github.com/slugger7/exorcist/internal/service/user.(*UserService).Create: could not determine if user '%v' exists\n%v", username, expectedError.Error())
	if _, err := us.Create(username, ""); err.Error() != expectedErrorMessage {
		t.Errorf("Unexpected error thrown\nExpected: %v\nGot: %v", expectedErrorMessage, err.Error())
	}
}

func Test_CreateUser_UserExistsTrue_ShouldReturnError(t *testing.T) {
	us, mr := setupMockUserService()

	expectedError := errors.New("user already exists")
	mr.mockUserRepo.mockModels[0] = &model.User{}

	username := ""

	if _, err := us.Create(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Unexpected error thrown\nExpected: %v\nGot: %v", expectedError.Error(), err.Error())
	}
}

func Test_CreateUser_UserExistsFalse_RepoCreateReturnsError_ShouldReturnError(t *testing.T) {
	us, mr := setupMockUserService()
	expectedError := errors.New("expected error")
	mr.mockUserRepo.mockErrors[1] = expectedError
	username := ""
	expectedErrorMessage := fmt.Sprintf("github.com/slugger7/exorcist/internal/service/user.(*UserService).Create: could not create a new user\n%v", expectedError.Error())

	if _, err := us.Create(username, ""); err.Error() != expectedErrorMessage {
		t.Errorf("Unexpected error thrown\nExpected: %v\nGot: %v", expectedErrorMessage, err.Error())
	}
}

func Test_CreateUser_UserExistsFalse_RepoCreatesUser_ShouldReturnUser(t *testing.T) {
	us, mr := setupMockUserService()
	username := "someUser"
	password := "somePassword"
	mr.mockUserRepo.mockModels[1] = &model.User{Username: username, Password: password}

	user, err := us.Create(username, "")
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Error("Created user is nil")
		return
	}
	if user.Username != username {
		t.Errorf("Unexpected user returned\nExpected: %v\nGot:%v", username, user.Username)
	}
	if user.Password != "" {
		t.Error("Password was not removed before returning")
	}
}

func Test_ValidateUser_RepoReturnsError_ShouldReturnError(t *testing.T) {
	us, mr := setupMockUserService()

	expecedError := errors.New("expected error")
	mr.mockUserRepo.mockErrors[0] = expecedError

	if _, err := us.Validate("", ""); err.Error() != expecedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expecedError.Error(), err.Error())
	}
}

func Test_ValidateUser_RepoReturnsNilUser_ShouldReturnError(t *testing.T) {
	us, _ := setupMockUserService()

	username := "someUsername"
	expectedError := fmt.Errorf("user with username %v does not exist", username)

	if _, err := us.Validate(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_ValidateUser_PasswordsDoNotMatch(t *testing.T) {
	us, mr := setupMockUserService()
	mr.mockUserRepo.mockModels[0] = &model.User{Password: ""}
	username := "someUsername"
	expectedError := fmt.Errorf("password for user %v did not match", username)

	if _, err := us.Validate(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_ValidateUser_PasswordsMatch_ShouldReturnUser(t *testing.T) {
	us, mr := setupMockUserService()

	password := "somePassword"
	username := "someUsername"
	mr.mockUserRepo.mockModels[0] = &model.User{Username: username, Password: hashPassword(password)}

	user, err := us.Validate(username, password)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Error("Expected user to be returned")
		return
	}
	if user.Username != username {
		t.Errorf("Unexpected user returned.\nExpected: %v\nGot: %v", username, user.Username)
	}
	if user.Password != "" {
		t.Errorf("Password was not cleared before returning")
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
