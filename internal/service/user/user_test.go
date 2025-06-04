package userService

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/assert"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
	mock_repository "github.com/slugger7/exorcist/internal/mock/repository"
	mock_userRepository "github.com/slugger7/exorcist/internal/mock/repository/user"
	"github.com/slugger7/exorcist/internal/models"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
	"go.uber.org/mock/gomock"
)

type testService struct {
	svc      *userService
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

	us := &userService{repo: mockRepo}
	return &testService{us, mockRepo, mockUserRepo}
}

func Test_UserExists_ErrorFromRepo(t *testing.T) {
	s := setup(t)

	expectedErr := fmt.Errorf("expected error")

	s.userRepo.EXPECT().
		GetUserByUsername(gomock.Any()).
		DoAndReturn(func(u string, columns ...postgres.Projection) (*model.User, error) {
			return nil, expectedErr
		}).
		Times(1)

	if _, err := s.svc.UserExists(""); err.Error() != expectedErr.Error() {
		t.Errorf("Expected: %v\nGot: %v", expectedErr.Error(), err.Error())
	}
}

func Test_UserExists_UserIsNil_ShouldReturnFalse(t *testing.T) {
	s := setup(t)

	s.userRepo.EXPECT().
		GetUserByUsername(gomock.Any()).
		DoAndReturn(func(u string, columns ...postgres.Projection) (*model.User, error) {
			return nil, nil
		}).
		Times(1)

	actual, err := s.svc.UserExists("")
	if err != nil {
		t.Fatal(err)
	}

	if actual {
		t.Error("Expected user to not exist")
	}
}

func Test_UserExists_UserIsDefined_ShouldReturnTrue(t *testing.T) {
	s := setup(t)

	s.userRepo.EXPECT().
		GetUserByUsername(gomock.Any()).
		DoAndReturn(func(u string, columns ...postgres.Projection) (*model.User, error) {
			return &model.User{}, nil
		}).
		Times(1)

	actual, err := s.svc.UserExists("")
	if err != nil {
		t.Fatal(err)
	}

	if !actual {
		t.Error("Expected user to exist")
	}
}

func Test_CreateUser_UserExistsRaisesError_ShouldReturnError(t *testing.T) {
	s := setup(t)

	s.userRepo.EXPECT().
		GetUserByUsername(gomock.Any()).
		DoAndReturn(func(u string, columns ...postgres.Projection) (*model.User, error) {
			return nil, fmt.Errorf("error")
		}).
		Times(1)

	username := "someUsername"

	user, err := s.svc.Create(username, "")
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
	s := setup(t)

	s.userRepo.EXPECT().
		GetUserByUsername(gomock.Any()).
		DoAndReturn(func(u string, columns ...postgres.Projection) (*model.User, error) {
			return &model.User{}, nil
		}).
		Times(1)
	username := "someUsername"

	user, err := s.svc.Create(username, "")
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
	s := setup(t)

	s.userRepo.EXPECT().
		CreateUser(gomock.Any()).
		DoAndReturn(func(user model.User) (*model.User, error) {
			return nil, fmt.Errorf("error")
		}).
		Times(1)

	s.userRepo.EXPECT().
		GetUserByUsername(gomock.Any()).
		DoAndReturn(func(u string, columns ...postgres.Projection) (*model.User, error) {
			return nil, nil
		}).
		Times(1)

	username := "someUsername"

	user, err := s.svc.Create(username, "")
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
	s := setup(t)

	username := "someUser"
	password := "somePassword"

	s.userRepo.EXPECT().
		GetUserByUsername(gomock.Any()).
		DoAndReturn(func(u string, columns ...postgres.Projection) (*model.User, error) {
			return nil, nil
		}).
		Times(1)
	s.userRepo.EXPECT().
		CreateUser(gomock.Any()).
		DoAndReturn(func(user model.User) (*model.User, error) {
			return &model.User{Username: username, Password: password}, nil
		}).
		Times(1)

	user, err := s.svc.Create(username, "")
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
	s := setup(t)

	expecedError := errors.New("expected error")

	s.userRepo.EXPECT().
		GetUserByUsername(gomock.Any(), gomock.Any()).
		DoAndReturn(func(username string, columns ...postgres.Projection) (*model.User, error) {
			return nil, expecedError
		}).
		Times(1)

	if _, err := s.svc.Validate("", ""); err.Error() != expecedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expecedError.Error(), err.Error())
	}
}

func Test_ValidateUser_RepoReturnsNilUser_ShouldReturnError(t *testing.T) {
	s := setup(t)

	username := "someUsername"
	expectedError := fmt.Errorf(ErrUserDoesNotExist, username)

	s.userRepo.EXPECT().
		GetUserByUsername(username, gomock.Any()).
		DoAndReturn(func(username string, columns ...postgres.Projection) (*model.User, error) {
			return nil, nil
		}).
		Times(1)

	if _, err := s.svc.Validate(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_ValidateUser_PasswordsDoNotMatch(t *testing.T) {
	s := setup(t)

	username := "someUsername"

	s.userRepo.EXPECT().
		GetUserByUsername(username, gomock.Any()).
		DoAndReturn(func(username string, columns ...postgres.Projection) (*model.User, error) {
			return &model.User{Password: ""}, nil
		}).
		Times(1)

	expectedError := fmt.Errorf(ErrUsersPasswordDidNotMatch, username)

	if _, err := s.svc.Validate(username, ""); err.Error() != expectedError.Error() {
		t.Errorf("Expected error: %v\nGot error: %v", expectedError.Error(), err.Error())
	}
}

func Test_ValidateUser_PasswordsMatch_ShouldReturnUser(t *testing.T) {
	s := setup(t)

	password := "somePassword"
	username := "someUsername"

	s.userRepo.EXPECT().
		GetUserByUsername(username, gomock.Any()).
		DoAndReturn(func(username string, columns ...postgres.Projection) (*model.User, error) {
			return &model.User{Username: username, Password: hashPassword(password)}, nil
		}).
		Times(1)

	user, err := s.svc.Validate(username, password)
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
