package userService

import (
	"errors"
	"fmt"
	"log"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateUser(username, password string) (*model.User, error)
	ValidateUser(username, password string) (*model.User, error)
}

type UserService struct {
	Env  *environment.EnvironmentVariables
	repo repository.IRepository
}

var userServiceInstance *UserService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) *UserService {
	if userServiceInstance == nil {
		userServiceInstance = &UserService{
			Env:  env,
			repo: repo,
		}

		log.Println("UserService instance created")
	}
	return userServiceInstance
}

func (us *UserService) UserExists(username string) (bool, error) {
	user, err := us.repo.UserRepo().GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func (us *UserService) CreateUser(username, password string) (*model.User, error) {
	log.Println("Creating user in service")
	userExists, err := us.UserExists(username)
	if err != nil {
		return nil, errors.Join(errors.New(fmt.Sprintf("could not determine if user '%v' exists", username)), err)
	}

	if userExists {
		return nil, errors.New("user already exists")
	}

	user := model.User{
		Username: username,
		Password: hashPassword(password),
	}

	newUser, err := us.repo.UserRepo().CreateUser(user)
	if err != nil {
		return nil, errors.Join(errors.New("could not create a new user"), err)
	}

	newUser.Password = ""

	return newUser, nil
}

func (us *UserService) ValidateUser(username, password string) (*model.User, error) {
	var users []struct {
		model.User
	}
	user, err := us.repo.UserRepo().
		GetUserByUsername(username, table.User.ID, table.User.Password)
	if err != nil {
		return nil, err
	}
	if len(users) > 1 {
		panic("Found more than one active user for a username")
	}

	if !compareHashedPassword(user.Password, password) {
		log.Printf("Password did not match hashed password in database for user %v", username)
		return nil, nil
	}
	user.Password = "" // do not want to return the hash of the password

	return user, nil
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	errs.CheckError(err)

	return string(hashedPassword)
}

func compareHashedPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
