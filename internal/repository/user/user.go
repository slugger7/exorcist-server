package userRepository

import (
	"database/sql"
	"log"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type UserStatement struct {
	postgres.Statement
	db *sql.DB
}

type IUserRepository interface {
	GetUserByUsernameAndPassword(username, password string) *UserStatement
	GetUserByUsername(username string, columns ...postgres.Projection) *UserStatement
	CreateUser(user model.User) *UserStatement
}

type UserRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var userRepositoryInstance *UserRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) IUserRepository {
	if userRepositoryInstance == nil {
		userRepositoryInstance = &UserRepository{
			db:  db,
			Env: env,
		}

		log.Println("User repository instance created")
	}

	return userRepositoryInstance
}

func (us *UserStatement) Query(destination interface{}) error {
	log.Println("Querying user statment")
	return us.Statement.Query(us.db, destination)
}

func (ur *UserRepository) GetUserByUsernameAndPassword(username, password string) *UserStatement {
	statement := table.User.SELECT(table.User.ID, table.User.Username).
		FROM(table.User).
		WHERE(table.User.Username.EQ(postgres.String(username)).
			AND(table.User.Password.EQ(postgres.String(password))).
			AND(table.User.Active.IS_TRUE()))

	util.DebugCheck(ur.Env, statement)
	return &UserStatement{statement, ur.db}
}

func (ur *UserRepository) GetUserByUsername(username string, columns ...postgres.Projection) *UserStatement {
	if len(columns) == 0 {
		columns = []postgres.Projection{table.User.Username}
	}
	statement := table.User.SELECT(columns[0], columns[1:]...).
		FROM(table.User).
		WHERE(table.User.Username.EQ(postgres.String(username)).
			AND(table.User.Active.IS_TRUE()))

	util.DebugCheck(ur.Env, statement)
	return &UserStatement{statement, ur.db}
}

func (ur *UserRepository) CreateUser(user model.User) *UserStatement {
	statement := table.User.INSERT(table.User.Username, table.User.Password).
		MODEL(user).
		RETURNING(table.User.ID, table.User.Username, table.User.Active, table.User.Created, table.User.Modified)

	util.DebugCheck(ur.Env, statement)
	return &UserStatement{statement, ur.db}
}
