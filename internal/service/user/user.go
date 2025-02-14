package userService

import (
	"errors"
	"fmt"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(username, password string) (*model.User, error)
	Validate(username, password string) (*model.User, error)
}

type UserService struct {
	Env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

var userServiceInstance *UserService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) *UserService {
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

func (us *UserService) Create(username, password string) (*model.User, error) {
	userExists, err := us.UserExists(username)
	if err != nil {
		return nil, errs.BuildError(err, "could not determine if user '%v' exists", username)
	}

	if userExists {
		return nil, errors.New("user already exists")
	}

	user := model.User{
		Username: username,
		Password: hashPassword(password),
	}

	newUser, err := us.repo.User().CreateUser(user)
	if err != nil {
		return nil, errs.BuildError(err, "could not create a new user")
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

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	errs.PanicError(err)

	return string(hashedPassword)
}

func compareHashedPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
