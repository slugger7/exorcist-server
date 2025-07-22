package libraryPathRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type libraryPathRepository struct {
	db  *sql.DB
	env *environment.EnvironmentVariables
	ctx context.Context
}

var libraryPathRepoInstance *libraryPathRepository

type LibraryPathStatement struct {
	postgres.Statement
	db  *sql.DB
	ctx context.Context
}

type LibraryPathRepository interface {
	Create(path string, libraryId uuid.UUID) (*model.LibraryPath, error)
	GetAll() ([]model.LibraryPath, error)
	GetById(id uuid.UUID) (*model.LibraryPath, error)
	GetByLibraryId(libraryId uuid.UUID) ([]model.LibraryPath, error)
	GetContainingPath(path string) ([]model.LibraryPath, error)
}

func (i *libraryPathRepository) GetContainingPath(path string) ([]model.LibraryPath, error) {
	libraryPath := table.LibraryPath
	stmnt := libraryPath.SELECT(libraryPath.ID, libraryPath.Path).
		WHERE(libraryPath.Path.LIKE(postgres.String(fmt.Sprintf("%v%%", path))).OR(
			postgres.BoolExp(postgres.Raw("#1 LIKE library_path.\"path\" || '%'", postgres.RawArgs{"#1": path})),
		))

	util.DebugCheck(i.env, stmnt)

	var libraryPaths []model.LibraryPath
	if err := stmnt.QueryContext(i.ctx, i.db, &libraryPaths); err != nil {
		return nil, errs.BuildError(err, "could not fetch library paths containing path: %v", path)
	}

	return libraryPaths, nil
}

func New(db *sql.DB, env *environment.EnvironmentVariables, context context.Context) LibraryPathRepository {
	if libraryPathRepoInstance != nil {
		return libraryPathRepoInstance
	}
	libraryPathRepoInstance = &libraryPathRepository{
		db:  db,
		env: env,
		ctx: context,
	}

	return libraryPathRepoInstance
}

func (lps LibraryPathStatement) Query(destination interface{}) error {
	return lps.Statement.QueryContext(lps.ctx, lps.db, destination)
}

func (lps *libraryPathRepository) Create(path string, libraryId uuid.UUID) (*model.LibraryPath, error) {
	var libraryPath struct{ model.LibraryPath }
	if err := lps.create(&model.LibraryPath{Path: path, LibraryID: libraryId}).Query(&libraryPath); err != nil {
		return nil, errs.BuildError(err, "could not create library path, with \npath: %v\nlibraryId: %v", path, libraryId)
	}
	return &libraryPath.LibraryPath, nil
}

func (lps *libraryPathRepository) GetAll() ([]model.LibraryPath, error) {
	var libraryPaths []model.LibraryPath
	if err := lps.getLibraryPathsSelect().Query(&libraryPaths); err != nil {
		return nil, errs.BuildError(err, "could not get library paths")
	}
	return libraryPaths, nil
}

func (lps *libraryPathRepository) GetByLibraryId(libraryId uuid.UUID) ([]model.LibraryPath, error) {
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

func (lps *libraryPathRepository) GetById(id uuid.UUID) (*model.LibraryPath, error) {
	var libraryPaths []struct{ model.LibraryPath }
	if err := lps.getByIdStatement(id).Query(&libraryPaths); err != nil {
		return nil, errs.BuildError(err, "could not get library path by id: %v", id)
	}

	if len(libraryPaths) != 1 {
		return nil, nil
	}

	return &libraryPaths[len(libraryPaths)-1].LibraryPath, nil
}
