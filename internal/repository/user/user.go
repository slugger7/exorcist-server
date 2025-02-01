package userRepository

import (
	"database/sql"
	"log"

	"github.com/go-jet/jet/v2/postgres"
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

func (ur *UserRepository) GetUserByUsernameAndPassword(username, password string) *UserStatement {
	statement := table.Users.SELECT(table.Users.ID, table.Users.Username).
		FROM(table.Users).
		WHERE(table.Users.Username.EQ(postgres.String(username)).
			AND(table.Users.Password.EQ(postgres.String(password))))

	util.DebugCheck(ur.Env, statement)
	return &UserStatement{statement, ur.db}
}

func (us *UserStatement) Query(destination interface{}) error {
	log.Println("Querying user statment")
	return us.Statement.Query(us.db, destination)
}
