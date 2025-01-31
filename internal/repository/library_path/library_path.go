package libraryPathRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type ILibraryPathRepository interface {
	GetLibraryPathsSelect() postgres.SelectStatement
	CreateLibraryPath(libraryId uuid.UUID, path string) postgres.InsertStatement
}

type LibraryPathRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var libraryPathRepoInstance *LibraryPathRepository

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

func (ds *LibraryPathRepository) GetLibraryPathsSelect() postgres.SelectStatement {
	selectQuery := table.LibraryPath.
		SELECT(table.LibraryPath.ID, table.LibraryPath.Path).
		FROM(table.LibraryPath)

	util.DebugCheck(ds.Env, selectQuery)
	return selectQuery
}

// TODO write test for function
func (ds *LibraryPathRepository) CreateLibraryPath(libraryId uuid.UUID, path string) postgres.InsertStatement {
	newLibPath := model.LibraryPath{
		LibraryID: libraryId,
		Path:      path,
	}

	insertStatement := table.LibraryPath.
		INSERT(
			table.LibraryPath.LibraryID,
			table.LibraryPath.Path,
		).
		MODEL(newLibPath).
		RETURNING(table.LibraryPath.ID, table.LibraryPath.Path)

	util.DebugCheck(ds.Env, insertStatement)

	return insertStatement
}
