package libraryPathRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
)

type ILibraryPathRepository interface {
	Create(*model.LibraryPath) (*model.LibraryPath, error)
}

type LibraryPathRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var libraryPathRepoInstance *LibraryPathRepository

type LibraryPathStatement struct {
	postgres.Statement
	db *sql.DB
}

func New(db *sql.DB, env *environment.EnvironmentVariables) ILibraryPathRepository {
	if libraryPathRepoInstance != nil {
		return libraryPathRepoInstance
	}
	libraryPathRepoInstance = &LibraryPathRepository{
		db:  db,
		Env: env,
	}

	return libraryPathRepoInstance
}

func (lps LibraryPathStatement) Query(destination interface{}) error {
	return lps.Statement.Query(lps.db, destination)
}

func (lps *LibraryPathRepository) Create(libraryPath *model.LibraryPath) (*model.LibraryPath, error) {
	panic("not implemented")
}
