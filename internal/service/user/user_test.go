package userService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	mock_repository "github.com/slugger7/exorcist/internal/mock/repository"
	mock_userRepository "github.com/slugger7/exorcist/internal/mock/repository/user"
	"github.com/slugger7/exorcist/internal/mocks/mrepository"
	"github.com/slugger7/exorcist/internal/models"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
	"go.uber.org/mock/gomock"
)

// Deprecated: use mockgen mocks instead
func setupOld() (*UserService, *mrepository.MockRepository) {
	mockRepo := mrepository.SetupMockRespository()
	us := &UserService{repo: mockRepo}
	return us, mockRepo
}

type testService struct {
	svc      *UserService
	repo     *mock_repository.MockIRepository
	userRepo *mock_userRepository.MockIUserRepository
}

func setup(t *testing.T) *testService {
	ctrl := gomock.NewController(t)

	mockRepo := mock_repository.NewMockIRepository(ctrl)
	mockUserRepo := mock_userRepository.NewMockIUserRepository(ctrl)

	mockRepo.EXPECT().
		User().
		DoAndReturn(func() userRepository.IUserRepository {
			return mockUserRepo
		}).
		AnyTimes()

	us := &UserService{repo: mockRepo}
	return &testService{us, mockRepo, mockUserRepo}
}

func Test_UserExists_ErrorFromRepo(t *testing.T) {
	us, mr := setupOld()
	expectedErr := errors.New("expected error")
	mr.MockUserRepo.MockError[0] = expectedErr

	if _, err := us.UserExists(""); err.Error() != expectedErr.Error() {
		t.Errorf("Expected: %v\nGot: %v", expectedErr.Error(), err.Error())
	}
}

func Test_UserExists_UserIsNil_ShouldReturnFalse(t *testing.T) {
	us, mr := setupOld()

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
	us, mr := setupOld()
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
	us, mr := setupOld()

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
	us, mr := setupOld()

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
	us, mr := setupOld()

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
	us, mr := setupOld()
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
	us, mr := setupOld()

	expecedError := errors.New("expected error")
	mr.MockUserRepo.MockError[0] = expecedError

	if _, err := us.Validate("", ""); err.Error() != expecedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expecedError.Error(), err.Error())
	}
}

func Test_ValidateUser_RepoReturnsNilUser_ShouldReturnError(t *testing.T) {
	us, _ := setupOld()

	username := "someUsername"
	expectedError := fmt.Errorf("user with username %v does not exist", username)

	if _, err := us.Validate(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_ValidateUser_PasswordsDoNotMatch(t *testing.T) {
	us, mr := setupOld()
	mr.MockUserRepo.MockModel[0] = &model.User{Password: ""}
	username := "someUsername"
	expectedError := fmt.Errorf("password for user %v did not match", username)

	if _, err := us.Validate(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_ValidateUser_PasswordsMatch_ShouldReturnUser(t *testing.T) {
	us, mr := setupOld()

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

func Test_UpdatePassword_GetUserReturnsError(t *testing.T) {
	s := setup(t)

	id, _ := uuid.NewRandom()

	s.userRepo.EXPECT().
		GetById(gomock.Eq(id)).
		DoAndReturn(func(uuid.UUID) (*model.User, error) {
			return nil, fmt.Errorf("some error")
		}).
		Times(1)

	err := s.svc.UpdatePassword(id, models.ResetPasswordModel{})

	assert.ErrorNotNil(t, err)
	assert.ErrorMessage(t, fmt.Sprintf(ErrGetById, id), err)
}

func Test_UpdatePassword_GetUserReturnsNilUser(t *testing.T) {
	s := setup(t)

	id, _ := uuid.NewRandom()

	s.userRepo.EXPECT().
		GetById(gomock.Eq(id)).
		DoAndReturn(func(uuid.UUID) (*model.User, error) {
			return nil, nil
		}).
		Times(1)

	err := s.svc.UpdatePassword(id, models.ResetPasswordModel{})

	assert.ErrorNotNil(t, err)
	assert.Error(t, fmt.Errorf(ErrUserNil, id), err)
}

func Test_UpdatePassword_PasswordsDoNotMatch(t *testing.T) {
	s := setup(t)

	id, _ := uuid.NewRandom()

	m := models.ResetPasswordModel{OldPassword: "someOldPassword"}
	u := model.User{Password: "thisPasswordWillNotMatch"}

	s.userRepo.EXPECT().
		GetById(gomock.Eq(id)).
		DoAndReturn(func(uuid.UUID) (*model.User, error) {
			return &u, nil
		}).
		Times(1)

	err := s.svc.UpdatePassword(id, m)

	assert.ErrorNotNil(t, err)
	assert.Error(t, fmt.Errorf(ErrNonMatchingPasswords, id), err)
}

func Test_UpdatePassword_RepoUpdateReturnsErr(t *testing.T) {
	s := setup(t)

	id, _ := uuid.NewRandom()

	m := models.ResetPasswordModel{OldPassword: "someOldPassword", NewPassword: "someNewPassword"}
	u := model.User{Password: hashPassword(m.OldPassword)}

	s.userRepo.EXPECT().
		GetById(gomock.Eq(id)).
		DoAndReturn(func(uuid.UUID) (*model.User, error) {
			return &u, nil
		}).
		Times(1)

	s.userRepo.EXPECT().
		UpdatePassword(gomock.Any()).
		DoAndReturn(func(*model.User) error {
			return fmt.Errorf("some error")
		})

	err := s.svc.UpdatePassword(id, m)

	assert.ErrorNotNil(t, err)
	assert.ErrorMessage(t, fmt.Sprintf(ErrUpdatingPassword, id), err)
}

func Test_UpdatePassword_Success(t *testing.T) {
	s := setup(t)

	id, _ := uuid.NewRandom()

	m := models.ResetPasswordModel{OldPassword: "someOldPassword", NewPassword: "someNewPassword"}
	u := model.User{Password: hashPassword(m.OldPassword)}

	s.userRepo.EXPECT().
		GetById(gomock.Eq(id)).
		DoAndReturn(func(uuid.UUID) (*model.User, error) {
			return &u, nil
		}).
		Times(1)

	s.userRepo.EXPECT().
		UpdatePassword(gomock.Any()).
		DoAndReturn(func(*model.User) error {
			return nil
		})

	err := s.svc.UpdatePassword(id, m)

	assert.ErrorNil(t, err)
}
