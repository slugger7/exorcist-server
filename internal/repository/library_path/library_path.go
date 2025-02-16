package libraryPathRepository

import (
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
	var libraryPath struct{ model.LibraryPath }
	if err := lps.getByIdStatement(id).Query(&libraryPath); err != nil {
		return nil, errs.BuildError(err, "could not get library path by id: %v", id)
	}

	return &libraryPath.LibraryPath, nil
}
