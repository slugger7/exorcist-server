package libraryRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type LibraryStatement struct {
	postgres.Statement
	db *sql.DB
}
type ILibraryRepository interface {
	CreateLibraryStatement(name string) ILibraryStatement
	GetLibraryByName(name string) ILibraryStatement
}

type ILibraryStatement interface {
	Query(destination interface{}) error
	Sql() string
}

type LibraryRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var libraryRepoInstance *LibraryRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) ILibraryRepository {
	if libraryRepoInstance != nil {
		return libraryRepoInstance
	}
	libraryRepoInstance = &LibraryRepository{
		db:  db,
		Env: env,
	}
	return libraryRepoInstance
}

func (ls *LibraryStatement) Query(destination interface{}) error {
	return ls.Statement.Query(ls.db, destination)
}

func (ls *LibraryStatement) Sql() string {
	sql, _ := ls.Statement.Sql()
	return sql
}

func (ls *LibraryRepository) CreateLibraryStatement(name string) ILibraryStatement {
	newLibrary := model.Library{
		Name: name,
	}

	insertStatement := table.Library.INSERT(table.Library.Name).
		MODEL(newLibrary).
		RETURNING(table.Library.ID)

	util.DebugCheck(ls.Env, insertStatement)

	return &LibraryStatement{insertStatement, ls.db}
}

func (i *LibraryRepository) GetLibraryByName(name string) ILibraryStatement {
	statement := table.Library.SELECT(table.Library.ID).
		FROM(table.Library).
		WHERE(table.Library.Name.EQ(postgres.String(name)))

	util.DebugCheck(i.Env, statement)
	return &LibraryStatement{statement, i.db}
}
