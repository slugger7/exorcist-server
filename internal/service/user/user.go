package userService

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	"github.com/slugger7/exorcist/internal/repository"
)

type UserService interface {
	Create(username, password string) (*model.User, error)
	Validate(username, password string) (*model.User, error)
	UpdatePassword(id uuid.UUID, model dto.ResetPasswordDTO) error
	AddMediaToFavourites(id, mediaId uuid.UUID) error
}

func (u *userService) AddMediaToFavourites(userId uuid.UUID, mediaId uuid.UUID) error {
	media, err := u.repo.Media().GetById(mediaId)
	if err != nil {
		return errs.BuildError(err, "could not get media by id: %v", mediaId.String())
	}

	if media == nil {
		return fmt.Errorf("media %v was nil and can't be added as favourite", mediaId.String())
	}

	favouriteMedia, err := u.repo.User().GetFavourite(userId, mediaId)
	if err != nil {
		return errs.BuildError(err, "")
	}

	if favouriteMedia != nil {
		u.logger.Warningf("favourite media already exists for user %v and media %v", userId.String(), mediaId.String())
		return nil
	}

	if err := u.repo.User().AddMediaToFavourites(userId, mediaId); err != nil {
		return errs.BuildError(err, "could not add media %v to favourites for %v", mediaId.String(), userId.String())
	}

	return nil
}

type userService struct {
	env    *environment.EnvironmentVariables
	repo   repository.IRepository
	logger logger.ILogger
}

var userServiceInstance *userService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) UserService {
	if userServiceInstance == nil {
		userServiceInstance = &userService{
			env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		userServiceInstance.logger.Info("UserService instance created")
	}
	return userServiceInstance
}

func (us *userService) UserExists(username string) (bool, error) {
	user, err := us.repo.User().GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

const ErrDeterminingUserExists = "could not determine if user '%v' exists"
const ErrUserExists = "user already exists"
const ErrCreatingUser = "could not create a new user"

func (us *userService) Create(username, password string) (*model.User, error) {
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

const (
	ErrUserDoesNotExist         string = "user with username %v does not exist"
	ErrUsersPasswordDidNotMatch string = "password for user %v did not match"
)

func (us *userService) Validate(username, password string) (*model.User, error) {
	user, err := us.repo.User().
		GetUserByUsername(username, table.User.ID, table.User.Password)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf(ErrUserDoesNotExist, username)
	}

	if !compareHashedPassword(user.Password, password) {
		return nil, fmt.Errorf(ErrUsersPasswordDidNotMatch, username)
	}
	user.Password = "" // do not want to return the hash of the password

	return user, nil
}

const (
	ErrGetById              string = "could not get user by id %v"
	ErrUserNil              string = "user with id %v does not exist"
	ErrNonMatchingPasswords string = "old password for user %v did not match"
	ErrUpdatingPassword     string = "could not update password for user %v"
)

func (us *userService) UpdatePassword(id uuid.UUID, m dto.ResetPasswordDTO) error {
	user, err := us.repo.User().GetById(id)
	if err != nil {
		return errs.BuildError(err, ErrGetById, id)
	}

	if user == nil {
		return fmt.Errorf(ErrUserNil, id)
	}

	if !compareHashedPassword(user.Password, m.OldPassword) {
		return fmt.Errorf(ErrNonMatchingPasswords, id)
	}

	user.Password = hashPassword(m.NewPassword)
	if err := us.repo.User().UpdatePassword(user); err != nil {
		return errs.BuildError(err, ErrUpdatingPassword, id)
	}

	return nil
}
