package userService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/mocks/mrepository"
)

func setup() (*UserService, *mrepository.MockRepository) {
	mockRepo := mrepository.SetupMockRespository()
	us := &UserService{repo: mockRepo}
	return us, mockRepo
}
func Test_UserExists_ErrorFromRepo(t *testing.T) {
	us, mr := setup()
	expectedErr := errors.New("expected error")
	mr.MockUserRepo.MockError[0] = expectedErr

	if _, err := us.UserExists(""); err.Error() != expectedErr.Error() {
		t.Errorf("Expected: %v\nGot: %v", expectedErr.Error(), err.Error())
	}
}

func Test_UserExists_UserIsNil_ShouldReturnFalse(t *testing.T) {
	us, mr := setup()

	mr.MockUserRepo.MockModel[0] = nil
	mr.MockUserRepo.MockError[0] = nil

	actual, err := us.UserExists("")
	if err != nil {
		t.Fatal(err)
	}

	if actual {
		t.Error("Expected user to not exist")
	}
}

func Test_UserExists_UserIsDefined_ShouldReturnTrue(t *testing.T) {
	us, mr := setup()
	mr.MockUserRepo.MockModel[0] = &model.User{}

	actual, err := us.UserExists("")
	if err != nil {
		t.Fatal(err)
	}

	if !actual {
		t.Error("Expected user to exist")
	}
}

func Test_CreateUser_UserExistsRaisesError_ShouldReturnError(t *testing.T) {
	us, mr := setup()

	mr.MockUserRepo.MockError[0] = errors.New("error")
	username := "someUsername"

	user, err := us.Create(username, "")
	if err == nil {
		t.Error("Expected an error but it was nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		expectedErr := fmt.Sprintf(ErrDeterminingUserExists, username)
		if e.Message() != expectedErr {
			t.Errorf("Expected error: %v\nGot error: %v", expectedErr, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got: %v", err)
	}

	if user != nil {
		t.Error("Error was raised but user was not nil")
	}
}

func Test_CreateUser_UserExistsTrue_ShouldReturnError(t *testing.T) {
	us, mr := setup()

	mr.MockUserRepo.MockModel[0] = &model.User{}

	username := "someUsername"

	user, err := us.Create(username, "")
	if err == nil {
		t.Error("Expected error but was nil")
	}
	if err.Error() != ErrUserExists {
		t.Errorf("Expected error: %v\nGot error: %v", ErrUserExists, err)
	}

	if user != nil {
		t.Error("Err was raised but user was not nil")
	}
}

func Test_CreateUser_UserExistsFalse_RepoCreateReturnsError_ShouldReturnError(t *testing.T) {
	us, mr := setup()

	mr.MockUserRepo.MockError[1] = fmt.Errorf("error")
	username := "someUsername"

	user, err := us.Create(username, "")
	if err == nil {
		t.Error("Expected an error but was nil")
	}
	var e errs.IError
	if errors.As(err, &e) {
		if e.Message() != ErrCreatingUser {
			t.Errorf("Expected error: %v\nGot error: %v", ErrCreatingUser, e.Message())
		}
	} else {
		t.Errorf("Expected specific error but got: %v", err)
	}

	if user != nil {
		t.Error("Expected an error and user to be nil but it was not")
	}
}

func Test_CreateUser_UserExistsFalse_RepoCreatesUser_ShouldReturnUser(t *testing.T) {
	us, mr := setup()
	username := "someUser"
	password := "somePassword"
	mr.MockUserRepo.MockModel[1] = &model.User{Username: username, Password: password}

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
	us, mr := setup()

	expecedError := errors.New("expected error")
	mr.MockUserRepo.MockError[0] = expecedError

	if _, err := us.Validate("", ""); err.Error() != expecedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expecedError.Error(), err.Error())
	}
}

func Test_ValidateUser_RepoReturnsNilUser_ShouldReturnError(t *testing.T) {
	us, _ := setup()

	username := "someUsername"
	expectedError := fmt.Errorf("user with username %v does not exist", username)

	if _, err := us.Validate(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_ValidateUser_PasswordsDoNotMatch(t *testing.T) {
	us, mr := setup()
	mr.MockUserRepo.MockModel[0] = &model.User{Password: ""}
	username := "someUsername"
	expectedError := fmt.Errorf("password for user %v did not match", username)

	if _, err := us.Validate(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_ValidateUser_PasswordsMatch_ShouldReturnUser(t *testing.T) {
	us, mr := setup()

	password := "somePassword"
	username := "someUsername"
	mr.MockUserRepo.MockModel[0] = &model.User{Username: username, Password: hashPassword(password)}

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
