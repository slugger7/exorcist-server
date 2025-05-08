package libraryRepository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type ILibraryRepository interface {
	Create(name string) (*model.Library, error)
	GetByName(name string) (*model.Library, error)
	GetAll() ([]model.Library, error)
	GetById(uuid.UUID) (*model.Library, error)
}

type LibraryRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
	ctx context.Context
}

var libraryRepoInstance *LibraryRepository

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) ILibraryRepository {
	if libraryRepoInstance != nil {
		return libraryRepoInstance
	}
	libraryRepoInstance = &LibraryRepository{
		db:  db,
		Env: env,
		ctx: context,
	}
	return libraryRepoInstance
}

func (ls *LibraryStatement) Query(destination interface{}) error {
	return ls.Statement.QueryContext(ls.ctx, ls.db, destination)
}

func (ls *LibraryStatement) Sql() string {
	sql, _ := ls.Statement.Sql()
	return sql
}

func (ls *LibraryRepository) Create(name string) (*model.Library, error) {
	var library struct{ model.Library }
	if err := ls.createLibraryStatement(name).Query(&library); err != nil {
		return nil, errs.BuildError(err, "error while creating library")
	}
	return &library.Library, nil
}

func (ls *LibraryRepository) GetByName(name string) (*model.Library, error) {
	var libraries []struct{ model.Library }
	if err := ls.getLibraryByNameStatement(name).Query(&libraries); err != nil {
		return nil, errs.BuildError(err, "could not get library by name '%v'", name)
	}
	var library *model.Library
	if len(libraries) > 0 {
		library = &libraries[len(libraries)-1].Library
	}
	return library, nil
}

func (ls *LibraryRepository) GetAll() ([]model.Library, error) {
	var libraries []struct{ model.Library }
	if err := ls.getLibrariesStatement().Query(&libraries); err != nil {
		return nil, errs.BuildError(err, "could not get libraries")
	}
	var libs []model.Library
	for _, v := range libraries {
		libs = append(libs, v.Library)
	}

	return libs, nil
}

func (ls *LibraryRepository) GetById(id uuid.UUID) (*model.Library, error) {
	var library struct{ model.Library }
	if err := ls.getById(id).Query(&library); err != nil {
		return nil, errs.BuildError(err, "could not get library by id: %v", id)
	}
	return &library.Library, nil
}
