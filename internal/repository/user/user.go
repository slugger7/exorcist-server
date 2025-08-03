package userRepository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/dto"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/models"
	"github.com/slugger7/exorcist/internal/repository/helpers"
)

type UserStatement struct {
	postgres.Statement
	db  *sql.DB
	ctx context.Context
}

type UserRepository interface {
	GetUserByUsernameAndPassword(username, password string) (*model.User, error)
	GetUserByUsername(username string, columns ...postgres.Projection) (*model.User, error)
	CreateUser(user model.User) (*model.User, error)
	GetById(id uuid.UUID) (*model.User, error)
	UpdatePassword(user *model.User) error
	GetFavourites(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
}

type userRepository struct {
	db  *sql.DB
	env *environment.EnvironmentVariables
	ctx context.Context
}

// GetFavourites implements UserRepository.
func (ur *userRepository) GetFavourites(id uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	relationFn := func(relationTable postgres.ReadableTable) postgres.ReadableTable {
		return relationTable.INNER_JOIN(
			table.FavouriteMedia,
			table.FavouriteMedia.MediaID.EQ(table.Media.ID),
		)
	}

	whereFn := func(whr postgres.BoolExpression) postgres.BoolExpression {
		return whr.AND(table.FavouriteMedia.UserID.EQ(postgres.UUID(id)))
	}

	mediaPage, err := helpers.QueryMediaOverview(id, search, relationFn, whereFn, ur.ctx, ur.db, ur.env)
	if err != nil {
		return nil, errs.BuildError(err, "could not query media overview for favourites")
	}

	return mediaPage, nil
}

var userRepositoryInstance *userRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) UserRepository {
	if userRepositoryInstance == nil {
		userRepositoryInstance = &userRepository{
			db:  db,
			env: env,
			ctx: context,
		}
	}

	return userRepositoryInstance
}

func (us *UserStatement) Query(destination interface{}) error {
	return us.Statement.QueryContext(us.ctx, us.db, destination)
}

func (us *UserStatement) Exec() (sql.Result, error) {
	return us.Statement.ExecContext(us.ctx, us.db)
}

func (ur *userRepository) CreateUser(user model.User) (*model.User, error) {
	var newUser struct{ model.User }
	if err := ur.createStatement(user).Query(&newUser); err != nil {
		return nil, errs.BuildError(err, "could not create user %v", user)
	}
	return &newUser.User, nil
}

func (ur *userRepository) GetUserByUsername(username string, columns ...postgres.Projection) (*model.User, error) {
	var users []struct{ model.User }
	if err := ur.getUserByUsernameStatement(username, columns...).Query(&users); err != nil {
		log.Println("something went wrong getting user by username")
		return nil, errs.BuildError(err, "could not get user by username '%v'", username)
	}
	var user *model.User
	if len(users) > 0 {
		user = &users[len(users)-1].User
	}
	return user, nil
}

func (ur *userRepository) GetUserByUsernameAndPassword(username, password string) (*model.User, error) {
	var users []struct{ model.User }
	if err := ur.getUserByUsernameAndPasswordStatement(username, password).Query(&users); err != nil {
		log.Println("something went wrong getting user by username")
		return nil, errs.BuildError(err, "could not get user by username '%v' and password", username)
	}
	var user *model.User
	if len(users) > 0 {
		user = &users[len(users)-1].User
	}
	return user, nil
}

func (ur *userRepository) GetById(id uuid.UUID) (*model.User, error) {
	var users []struct{ model.User }
	if err := ur.getByIdStatement(id).Query(&users); err != nil {
		return nil, errs.BuildError(err, "could not get user by id: %v", id)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no user found by id: %v", id)
	}

	return &users[len(users)-1].User, nil
}

func (ur *userRepository) UpdatePassword(user *model.User) error {
	if _, err := ur.updatePasswordStatement(user).Exec(); err != nil {
		return errs.BuildError(err, "could not update user password: %v", user.ID)
	}

	return nil
}
