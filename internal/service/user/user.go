package userService

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository"
)

type IUserService interface {
	Create(username, password string) (*model.User, error)
	Validate(username, password string) (*model.User, error)
	UpdatePassword(id uuid.UUID, model models.ResetPasswordModel) error
}

type UserService struct {
	Env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

var userServiceInstance *UserService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) IUserService {
	if userServiceInstance == nil {
		userServiceInstance = &UserService{
			Env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		userServiceInstance.logger.Info("UserService instance created")
	}
	return userServiceInstance
}

func (us *UserService) UserExists(username string) (bool, error) {
	user, err := us.repo.User().GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

const ErrDeterminingUserExists = "could not determine if user '%v' exists"
const ErrUserExists = "user already exists"
const ErrCreatingUser = "could not create a new user"

func (us *UserService) Create(username, password string) (*model.User, error) {
	userExists, err := us.UserExists(username)
	if err != nil {
		return nil, errs.BuildError(err, ErrDeterminingUserExists, username)
	}

	if userExists {
		return nil, errors.New(ErrUserExists)
	}

	user := model.User{
		Username: username,
		Password: hashPassword(password),
	}

	newUser, err := us.repo.User().CreateUser(user)
	if err != nil {
		return nil, errs.BuildError(err, ErrCreatingUser)
	}

	newUser.Password = ""

	return newUser, nil
}

func (us *UserService) Validate(username, password string) (*model.User, error) {
	user, err := us.repo.User().
		GetUserByUsername(username, table.User.ID, table.User.Password)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user with username %v does not exist", username)
	}

	if !compareHashedPassword(user.Password, password) {
		return nil, fmt.Errorf("password for user %v did not match", username)
	}
	user.Password = "" // do not want to return the hash of the password

	return user, nil
}

const ErrGetById string = "could not get user by id %v"

func (us *UserService) UpdatePassword(id uuid.UUID, m models.ResetPasswordModel) error {
	user, err := us.repo.User().GetById(id)
	if err != nil {
		return errs.BuildError(err, ErrGetById, id)
	}

	if user == nil {
		return fmt.Errorf("user with id %v does not exist", id)
	}

	if !compareHashedPassword(user.Password, m.OldPassword) {
		return fmt.Errorf("old password for user %v did not match", id)
	}

	user.Password = hashPassword(m.NewPassword)
	if err := us.repo.User().UpdatePassword(user); err != nil {
		return errs.BuildError(err, "could not update password for user %v", id)
	}

	return nil
}
