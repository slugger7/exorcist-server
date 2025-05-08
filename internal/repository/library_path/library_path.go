package libraryPathRepository

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type ILibraryPathRepository interface {
	Create(path string, libraryId uuid.UUID) (*model.LibraryPath, error)
	GetAll() ([]model.LibraryPath, error)
	GetById(id uuid.UUID) (*model.LibraryPath, error)
	GetByLibraryId(libraryId uuid.UUID) ([]model.LibraryPath, error)
}

type LibraryPathRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
	ctx context.Context
}

var libraryPathRepoInstance *LibraryPathRepository

type LibraryPathStatement struct {
	postgres.Statement
	db  *sql.DB
	ctx context.Context
}

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) ILibraryPathRepository {
	if libraryPathRepoInstance != nil {
		return libraryPathRepoInstance
	}
	libraryPathRepoInstance = &LibraryPathRepository{
		db:  db,
		Env: env,
		ctx: context,
	}

	return libraryPathRepoInstance
}

func (lps LibraryPathStatement) Query(destination interface{}) error {
	return lps.Statement.QueryContext(lps.ctx, lps.db, destination)
}

func (lps *LibraryPathRepository) Create(path string, libraryId uuid.UUID) (*model.LibraryPath, error) {
	var libraryPath struct{ model.LibraryPath }
	if err := lps.create(&model.LibraryPath{Path: path, LibraryID: libraryId}).Query(&libraryPath); err != nil {
		return nil, errs.BuildError(err, "could not create library path, with \npath: %v\nlibraryId: %v", path, libraryId)
	}
	return &libraryPath.LibraryPath, nil
}

func (lps *LibraryPathRepository) GetAll() ([]model.LibraryPath, error) {
	var libraryPaths []struct{ model.LibraryPath }
	if err := lps.getLibraryPathsSelect().Query(&libraryPaths); err != nil {
		return nil, errs.BuildError(err, "could not get library paths")
	}
	libPathModels := make([]model.LibraryPath, len(libraryPaths))
	for _, l := range libraryPaths {
		libPathModels = append(libPathModels, l.LibraryPath)
	}
	return libPathModels, nil
}

func (lps *LibraryPathRepository) GetByLibraryId(libraryId uuid.UUID) ([]model.LibraryPath, error) {
	var libraryPaths []struct{ model.LibraryPath }
	if err := lps.getByLibraryIdStatement(libraryId).Query(&libraryPaths); err != nil {
		return nil, errs.BuildError(err, "could not get library paths for library %v", libraryId)
	}

	libPathModels := []model.LibraryPath{}
	for _, l := range libraryPaths {
		libPathModels = append(libPathModels, l.LibraryPath)
	}

	return libPathModels, nil
}

func (lps *LibraryPathRepository) GetById(id uuid.UUID) (*model.LibraryPath, error) {
	var libraryPaths []struct{ model.LibraryPath }
	if err := lps.getByIdStatement(id).Query(&libraryPaths); err != nil {
		return nil, errs.BuildError(err, "could not get library path by id: %v", id)
	}

	if len(libraryPaths) != 1 {
		return nil, nil
	}

	return &libraryPaths[len(libraryPaths)-1].LibraryPath, nil
}
