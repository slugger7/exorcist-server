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
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/repository/util"
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
	AddMediaToFavourites(userId uuid.UUID, mediaId uuid.UUID) error
	GetFavourite(id, mediaId uuid.UUID) (*model.FavouriteMedia, error)
	RemoveFavourite(id, mediaId uuid.UUID) error
}

type userRepository struct {
	db  *sql.DB
	env *environment.EnvironmentVariables
	ctx context.Context
}

func (u *userRepository) RemoveFavourite(id, mediaId uuid.UUID) error {
	statement := table.FavouriteMedia.DELETE().
		WHERE(table.FavouriteMedia.UserID.EQ(postgres.UUID(id)).
			AND(table.FavouriteMedia.MediaID.EQ(postgres.UUID(mediaId))))

	util.DebugCheck(u.env, statement)

	_, err := statement.ExecContext(u.ctx, u.db)
	if err != nil {
		return errs.BuildError(err, "could not remove favourite media for user %v and media %v", id.String(), mediaId.String())
	}

	return nil
}

func (u *userRepository) AddMediaToFavourites(userId uuid.UUID, mediaId uuid.UUID) error {
	statement := table.FavouriteMedia.INSERT(table.FavouriteMedia.UserID, table.FavouriteMedia.MediaID).
		MODEL(model.FavouriteMedia{UserID: userId, MediaID: mediaId})

	util.DebugCheck(u.env, statement)

	if _, err := statement.ExecContext(u.ctx, u.db); err != nil {
		return errs.BuildError(err, "could not insert favourite media record for media %v and user %v", mediaId.String(), userId.String())
	}

	return nil
}

// GetFavourite implements UserRepository.
func (ur *userRepository) GetFavourite(id uuid.UUID, mediaId uuid.UUID) (*model.FavouriteMedia, error) {
	statement := table.FavouriteMedia.SELECT(table.FavouriteMedia.AllColumns).
		WHERE(table.FavouriteMedia.UserID.EQ(postgres.UUID(id)).
			AND(table.FavouriteMedia.MediaID.EQ(postgres.UUID(mediaId))))

	util.DebugCheck(ur.env, statement)

	var favourites []model.FavouriteMedia
	if err := statement.QueryContext(ur.ctx, ur.db, &favourites); err != nil {
		return nil, errs.BuildError(err, "could not fetch favourite media entity for media %v and user id %v", mediaId.String(), id.String())
	}

	if len(favourites) == 0 {
		return nil, nil
	}

	return &favourites[0], nil
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
